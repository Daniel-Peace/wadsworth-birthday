import type {
  SlashCommandBuilder,
  ChatInputCommandInteraction,
} from "discord.js";

export interface Command {
  data: SlashCommandBuilder;
  execute: (interaction: ChatInputCommandInteraction) => Promise<void>;
}

export interface GuildUserPair {
  guildId: string;
  userId: string;
}

export interface BirthdayDocument {
  guildUserPair: GuildUserPair;
  day: number;
  month: number;
}

export interface DatabaseResponse {
  Status: number;
  Description: string;
  Data: string;
}
