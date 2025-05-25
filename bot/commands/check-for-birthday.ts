import { SlashCommandBuilder } from "discord.js";
import type {
  GuildUserPair,
  BirthdayDocument,
  DatabaseResponse,
} from "../types/Types";

export const prod_ready = false;

export const data = new SlashCommandBuilder()
  .setName("check-for-birthday")
  .setDescription("Wadsworth checks if your birthday is saved in the database");

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

  get_birthday(guildUserPair).then((dbResponse) => {
    switch (dbResponse.Status) {
      case 0:
        if (dbResponse.Data === "") {
          interaction.reply(
            `I couldn't find your birthday in my database. You can add using the \`add-birthday\` command if you would like though.`,
          );
        } else {
          const birthdayDocument = JSON.parse(
            dbResponse.Data,
          ) as BirthdayDocument;

          interaction.reply(
            `I have ${birthdayDocument.Month}/${birthdayDocument.Day} as your birthday in my databse.`,
          );
        }

        break;
      default:
        interaction.reply(
          `Bummer! It Looks like something went wrong on my end. Maybe try again a bit later. Sorry about the inconvenience!`,
        );
        break;
    }
  });
}

async function get_birthday(doc: GuildUserPair): Promise<DatabaseResponse> {
  const params = new URLSearchParams({
    GuildId: doc.GuildId,
    UserId: doc.UserId,
  }).toString();
  console.log(doc.toString());
  const response = await fetch(
    `http://localhost:9000/check-for-bday?${params.toString()}`,
    {
      method: "GET",
    },
  );

  console.log(response);

  const data = await response.json();
  return data as DatabaseResponse;
}
