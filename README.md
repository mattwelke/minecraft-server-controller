# minecraft-server-controller

A small Go program that acts as a web server used to start and stop two Minecraft servers (one called "minecraft" and the other called "minecraft-private") Minecraft server.

File `main.go` can be modified if you want to use this software for your own purposes (e.g., only working with one server, adding more servers).

## Server setup

1. Ensure the server has a non-root user created to use to run this. The user must have sudo permission.

1. Add lines to the `sudoers` file to allow the user to start and stop the Minecraft server. This allows this server software to run the `systemctl start` and `systemctl stop` commands without having to use a password.

   Example:

   ```
   your_user ALL=(ALL) NOPASSWD: /bin/systemctl start minecraft
   your_user ALL=(ALL) NOPASSWD: /bin/systemctl stop minecraft
   ```

1. Set up a service to run this server software. This will work similar to how you set up a service to run your Minecraft server software. Ensure your service file starts this server software with the env var `SHUTDOWN_PASSWORD` set to a password the users will specify as they start and stop the server.

## Usage

1. Open `<host>:8080` in the web browser. Follow the on screen instructions.

   ![image](https://github.com/mattwelke/minecraft-server-controller/assets/7719209/ec1223a1-554d-454b-a3e0-4cdaea70f528)
