import { SlashCommandBuilder } from "discord.js";
import type { DatabaseResponse, GuildUserPair } from "../types/Types.ts";

export const prod_ready = false;

export const data = new SlashCommandBuilder()
  .setName("delete-birthday")
  .setDescription(
    "I will delete your birthday from my databse and will no longer wish you a happy birthday",
  );

export async function execute(interaction: any) {
  // get data
  const guildId = interaction.guildId;
  const userId = interaction.user.username;

  // create doc to send
  const guildUserPair: GuildUserPair = {
    GuildId: guildId,
    UserId: userId,
  };

  console.log(guildUserPair);

  delete_birthday(guildUserPair).then((dbResponse) => {
    switch (dbResponse.Status) {
      case 0:
        interaction.reply(
          `I deleted your birthday from my database and will no longer wish you a happy birthday. If you ever change your mind feel free to use the \`add-birthday\` commanmd to save it again.`,
        );
        break;
      case 1:
        interaction.reply(
          `Looks like I never had your birthday saved. If you ever do want to add your birthday feel free to use the \`add-birthday\` commanmd.`,
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

async function delete_birthday(doc: GuildUserPair): Promise<DatabaseResponse> {
  const response = await fetch("http://localhost:9000/delete-bday", {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(doc),
  });

  console.log(response);

  const data = await response.json();
  return data as DatabaseResponse;
}
