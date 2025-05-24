import { Client, Collection, Events, GatewayIntentBits } from "discord.js";
import { token } from "./botConfig.json";
import fs from "fs";
import path from "path";
import type { Command } from "./types/Types";

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

client.login(token);
