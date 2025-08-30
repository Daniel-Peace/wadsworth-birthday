import {
  SlashCommandBuilder,
  ChannelType,
  GuildMember,
  PermissionsBitField,
} from "discord.js";
import type { GuildConfig, DatabaseResponse } from "../types/Types";

export const prod_ready = false;

export const data = new SlashCommandBuilder()
  .setName("update-birthday-channel")
  .setDescription(
    "Wadsworth udpates the channel that he uses to wish people a happy birthday",
  )
  .addChannelOption((option) =>
    option
      .setName("channel")
      .setDescription("the channel you would Wadsworth to use")
      .addChannelTypes(ChannelType.GuildText)
      .setRequired(true),
  );

export async function execute(interaction: any) {
  const member = interaction.member;

  let isAdmin: boolean;
  if (member && "permissions" in member) {
    if (member.permissions.has(PermissionsBitField.Flags.Administrator)) {
      console.log("User is an admin");
      isAdmin = true;
    } else {
      console.log("User is not an admin");
      isAdmin = false;
    }
  } else {
    console.log("User is not an admin");
    isAdmin = false;
  }

  if (isAdmin) {
    const guildId = interaction.guildId;
    const channel = interaction.options.getChannel("channel");

    console.log(channel.id);

    // create doc to send
    const guildConfig: GuildConfig = {
      GuildId: guildId,
      RoleId: "",
      ChannelId: channel.id,
    };

    console.log(guildConfig);

    patch_config(guildConfig).then((dbResponse) => {
      switch (dbResponse.Status) {
        case 0:
          interaction.reply(
            `I updated your server config to use the channel \`${channel.name}\``,
          );
          break;
        default:
          interaction.reply(
            `Bummer! It Looks like something went wrong on my end. Maybe try again a bit later. Sorry about the inconvenience!`,
          );
          break;
      }
    });
  } else {
    await interaction.reply("Only admins may use this command");
  }
}

async function patch_config(doc: GuildConfig): Promise<DatabaseResponse> {
  const response = await fetch("http://localhost:9000/update-config", {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(doc),
  });

  console.log(response);

  const data = await response.json();
  return data as DatabaseResponse;
}
