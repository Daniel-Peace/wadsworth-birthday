import { SlashCommandBuilder, PermissionsBitField } from "discord.js";
import type { GuildConfig, DatabaseResponse } from "../types/Types";

export const prod_ready = false;

export const data = new SlashCommandBuilder()
  .setName("update-birthday-role")
  .setDescription(
    "Wadsworth udpates the role assigned to people with an active birthday",
  )
  .addRoleOption((option) =>
    option
      .setName("role")
      .setDescription("the role you would like to assign")
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
      console.log("User is an admin");
      isAdmin = false;
    }
  } else {
    console.log("User is an admin");
    isAdmin = false;
  }

  if (isAdmin) {
    const guildId = interaction.guildId;
    const role = interaction.options.getRole("role");

    // create doc to send
    const guildConfig: GuildConfig = {
      GuildId: guildId,
      RoleId: "",
      ChannelId: role.id,
    };

    console.log(guildConfig);

    patch_config(guildConfig).then((dbResponse) => {
      switch (dbResponse.Status) {
        case 0:
          interaction.reply(
            `I updated your server config to use the role \`${role.name}\``,
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
