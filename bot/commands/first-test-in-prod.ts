import { SlashCommandBuilder } from "discord.js";

export const prod_ready = true;

export const data = new SlashCommandBuilder()
  .setName("test-in-prod")
  .setDescription("Testing in prod!");

export async function execute(interaction: any) {
  await interaction.reply("Never fails! 😎");
}
