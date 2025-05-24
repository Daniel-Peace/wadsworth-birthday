import { SlashCommandBuilder } from "discord.js";
import type { GuildUserPair, BirthdayDocument } from "../types/Types";

export const prod_ready = false;

export const data = new SlashCommandBuilder()
  .setName("add-birthday")
  .setDescription(
    "Saves your birthday to my database so I can wish you a happy birthday.",
  )
  .addIntegerOption((option) =>
    option
      .setName("month")
      .setDescription("The month you were born in")
      .setRequired(true)
      .setMinValue(1)
      .setMaxValue(12),
  )

  .addIntegerOption((option) =>
    option
      .setName("day")
      .setDescription("The day you were born on")
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
    guildId: guildId,
    userId: userId,
  };
  const birthdayDocument: BirthdayDocument = {
    guildUserPair: guildUserPair,
    month: month,
    day: day,
  };

  console.log(birthdayDocument);

  send_post(birthdayDocument);

  await interaction.reply(`You entered the birthday ${month}/${day}`);
}

async function send_post(doc: BirthdayDocument) {
  fetch("http://localhost:9000/insert-bday", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(doc),
  })
    .then((response) => response.json())
    .then((json) => console.log(json));
}
