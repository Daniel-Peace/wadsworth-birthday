import { REST, Routes, SlashCommandBuilder } from "discord.js";
import fs from "fs";
import path from "path";
import { token, clientId, testGuildId } from "./botConfig.json";
import { exit } from "process";
import type { Command } from "./types/Types";
import type { RESTPostAPIApplicationCommandsJSONBody } from "discord.js";

const args = process.argv.slice(2);
const option = validate_args(args);

// Iterating through the arguments
args.forEach((arg, index) => {
  console.log(`Argument ${index + 1}: ${arg}`);
});

const rest = new REST({ version: "10" }).setToken(token);

const test_commands: RESTPostAPIApplicationCommandsJSONBody[] = [];
const prod_commands: RESTPostAPIApplicationCommandsJSONBody[] = [];

// getting all TS and JS files in the command directory
const commandsPath = path.join(__dirname, "commands");
const commandFiles = fs
  .readdirSync(commandsPath)
  .filter((file) => file.endsWith(".ts") || file.endsWith(".js"));

//looping through commands and adding them to their respective arrays
for (const file of commandFiles) {
  const filePath = path.join(commandsPath, file);
  const current_command = await import(filePath);
  if ("data" in current_command && "execute" in current_command) {
    if ("prod_ready" in current_command && current_command.prod_ready) {
      prod_commands.push(current_command.data.toJSON());
    }
    test_commands.push(current_command.data.toJSON());
  }
}

console.log("------------------------");
console.log(`TEST COMMANDS:`);
console.log("------------------------");
for (const command of test_commands) {
  console.log(`data:\n---\n${command}\n---`);
}
console.log("------------------------");
console.log(`PROD COMMANDS:`);
console.log("------------------------");
for (const command of prod_commands) {
  console.log(`data:\n---\n${command}\n---`);
}
console.log("------------------------");

switch (option) {
  case "all":
    register_test_commands();
    register_prod_commands();
    break;
  case "test":
    register_test_commands();
    break;
  case "prod":
    register_prod_commands();
    break;
  default:
    console.log("No arg found for how to deploy. Exiting...");
    exit(1);
}

function validate_args(args: string[]): string {
  if (args.length !== 1) {
    console.log(
      `[ERROR] - incorrect number of args [required: 1 | found: ${args.length}]`,
    );
    console.log("options:\n- all\n- test\n- prod");
    console.log("\t$ bun run register-commands.ts <option>");
    exit(1);
  } else if (args[0] === undefined) {
    console.log("[ERROR] - arg undefined");
    console.log("options:\n- all\n- test\n- prod");
    console.log("\t$ bun run register-commands.ts <option>");
    exit(1);
  } else {
    return args[0];
  }
}

async function register_prod_commands() {
  try {
    console.log("Registering commands for the production servers...");
    await rest.put(Routes.applicationCommands(clientId), {
      body: prod_commands,
    });
    console.log(`Registered ${prod_commands.length} commnds to prod servers`);
  } catch (error) {
    console.error(error);
  }
}

async function register_test_commands() {
  try {
    console.log("Registering commands for the test server...");
    await rest.put(Routes.applicationGuildCommands(clientId, testGuildId), {
      body: test_commands,
    });
    console.log(
      `Registered ${test_commands.length} commnds to the test server`,
    );
  } catch (error) {
    console.error(error);
  }
}
