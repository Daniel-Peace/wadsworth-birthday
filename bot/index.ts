import {
  Client,
  Collection,
  Events,
  GatewayIntentBits,
  ChannelType,
} from "discord.js";
import { token } from "./botConfig.json";
import fs from "fs";
import path from "path";
import type {
  Command,
  BirthdayDocument,
  DatabaseResponse,
} from "./types/Types";
import { sys } from "typescript";

// creating collection for commands
const client = new Client({ intents: [GatewayIntentBits.Guilds] }) as Client & {
  commands: Collection<string, Command>;
};
client.commands = new Collection();

// creating path to the commands directory
const commandsPath = path.join(__dirname, "commands");
const commandFiles = fs
  .readdirSync(commandsPath)
  .filter((file) => file.endsWith(".ts") || file.endsWith(".js"));

// importing each command
for (const file of commandFiles) {
  const command = await import(path.join(commandsPath, file));
  client.commands.set(command.data.name, command);
}

client.once(Events.ClientReady, () => {
  console.log(`Logged in as ${client.user!.tag}`);
});

client.on(Events.InteractionCreate, async (interaction) => {
  if (!interaction.isChatInputCommand()) return;

  const command = client.commands.get(interaction.commandName);
  if (!command) return;

  try {
    await command.execute(interaction);
  } catch (err) {
    console.error(err);
    await interaction.reply({
      content: "There was an error!",
      ephemeral: true,
    });
  }
});

scheduleNextDailyRun();

wishPeopleAHappyBirthday();

client.login(token);

function scheduleNextDailyRun() {
  const now = new Date();
  const nextRun = new Date();

  nextRun.setHours(8, 0, 0, 0);

  if (now >= nextRun) {
    nextRun.setDate(nextRun.getDate() + 1);
  }

  const msUntilNextRun = nextRun.getTime() - now.getTime();
  console.log("Next run in", msUntilNextRun / 1000, "seconds");

  setTimeout(async () => {
    await wishPeopleAHappyBirthday();
    scheduleNextDailyRun();
  }, msUntilNextRun);
}

async function wishPeopleAHappyBirthday() {
  const birthdays = await getActiveBirthdays();
  for (const birthday of birthdays) {
    const guild = await client.guilds.fetch(birthday.GuildUserPair.GuildId);
    const systemChannelId = guild.systemChannelId;
    if (systemChannelId) {
      const systemChannel = await guild.channels.fetch(systemChannelId);
      if (systemChannel && systemChannel.isTextBased()) {
        await systemChannel.send(
          `Happy Birthday ${birthday.GuildUserPair.UserId}`,
        );
        console.log("Message sent.");
      } else {
        console.log(
          "System channel is not text-based or could not be fetched.",
        );
      }
    } else {
      console.log("No system channel set for this guild.");
    }

    console.log(birthday);
  }
}

async function getActiveBirthdays(): Promise<BirthdayDocument[]> {
  console.log("getting active birthdays...");
  const response = await fetch("http://localhost:9000/get-active-bday", {
    method: "GET",
  });

  console.log(response);

  const data = await response.json();
  const databaseResponse = data as DatabaseResponse;

  if (databaseResponse.Data !== "") {
    return JSON.parse(databaseResponse.Data) as BirthdayDocument[];
  } else {
    return [];
  }
}
