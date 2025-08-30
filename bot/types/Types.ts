import type {
  SlashCommandBuilder,
  ChatInputCommandInteraction,
} from "discord.js";

export interface Command {
  data: SlashCommandBuilder;
  execute: (interaction: ChatInputCommandInteraction) => Promise<void>;
}

export interface GuildUserPair {
  GuildId: string;
  UserId: string;
}

export interface BirthdayDocument {
  GuildUserPair: GuildUserPair;
  Day: number;
  Month: number;
}

export interface DatabaseResponse {
  Status: number;
  Description: string;
  Data: string;
}

export interface GuildConfig {
  GuildId: string;
  RoleId: string;
  ChannelId: string;
}
