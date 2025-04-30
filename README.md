# vtomluml

# Vector.dev Config to PlantUML Diagram

This repository provides a tool to convert a Vector.dev configuration TOML file into a PlantUML diagram.

## Steps to Use

1. **Install Dependencies**
   ```sh
   go mod vendor
   ```

2. **Run the Tool**
   ```sh
   go run main.go
   ```

3. **Input File**
   - Store the content of your Vector.dev TOML configuration file in `data.toml`.

4. **Output File**
   - The PlantUML code will be generated in `diagram.puml`.
   - Upload the file in plantuml server or vscode plugin to generate the diagram.

## Example

1. Place your Vector.dev configuration in `data.toml`.
2. Execute the following commands:
   ```sh
   go mod vendor
   go run main.go
   ```
3. Check the generated `diagram.puml` for the PlantUML diagram.
