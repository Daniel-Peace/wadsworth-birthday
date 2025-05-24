import { SlashCommandBuilder } from "discord.js";

export const prod_ready = true;

export const data = new SlashCommandBuilder()
  .setName("second-test-in-prod")
  .setDescription("Testing in prod... again!");

export async function execute(interaction: any) {
  await interaction.reply("Never fails... again! 😎");
}
