import { SlashCommandBuilder } from "discord.js";
import type {
  GuildUserPair,
  BirthdayDocument,
  DatabaseResponse,
} from "../types/Types";

export const prod_ready = false;

export const data = new SlashCommandBuilder()
  .setName("update-birthday")
  .setDescription("Wadsworth udpates your birthday in the database")
  .addIntegerOption((option) =>
    option
      .setName("month")
      .setDescription("The month you were born in.")
      .setRequired(true)
      .setMinValue(1)
      .setMaxValue(12),
  )

  .addIntegerOption((option) =>
    option
      .setName("day")
      .setDescription("The day you were born on.")
      .setRequired(true)
      .setMinValue(1)
      .setMaxValue(31),
  );

export async function execute(interaction: any) {
  // get data
  const guildId = interaction.guildId;
  const userId = interaction.user.username;
  const month = interaction.options.getInteger("month");
  const day = interaction.options.getInteger("day");

  // create doc to send
  const guildUserPair: GuildUserPair = {
    GuildId: guildId,
    UserId: userId,
  };
  const birthdayDocument: BirthdayDocument = {
    GuildUserPair: guildUserPair,
    Month: month,
    Day: day,
  };

  console.log(birthdayDocument);

  put_birthday(birthdayDocument).then((dbResponse) => {
    switch (dbResponse.Status) {
      case 0:
        interaction.reply(`I updated your birthday to ${month}/${day}.`);
        break;
      case 1:
        interaction.reply(
          `Looks like there is no birthday to update, but you can add your birthday using the \`add-birthday\` command.`,
        );
        break;
      default:
        interaction.reply(
          `Bummer! It Looks like something went wrong on my end. Maybe try again a bit later. Sorry about the inconvenience!`,
        );
        break;
    }
  });
}

async function put_birthday(doc: BirthdayDocument): Promise<DatabaseResponse> {
  const response = await fetch("http://localhost:9000/update-bday", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(doc),
  });

  console.log(response);

  const data = await response.json();
  return data as DatabaseResponse;
}
