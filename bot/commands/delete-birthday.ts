import { SlashCommandBuilder } from "discord.js";

export const prod_ready = false;

export const data = new SlashCommandBuilder()
  .setName("delete-birthday")
  .setDescription(
    "I will delete your birthday from my databse and will no longer wish you a happy birthday",
  );

export async function execute(interaction: any) {
  await interaction.reply("You birthday has been successfully deleted!");
}
