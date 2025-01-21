# Steam Game Prices Notifier

## System Diagram

![steam_game_prices_notifier drawio](/docs/steam_game_prices_notifier.drawio.png)

## Application Description

Steam Game Prices Notifier tells you the best timing to buy video games on Steam.

- This app is synchronized to your Steam wishlist.
- Notion DB is used to store information of the current and lowest prices of games.
- If the current prices of games are cheaper than or equal to their lowest prices recorded in the Notion DB, the app automatically notifies you prices of those games.
- This app runs at 18:00 pm (JST) every day.

## How to Set up the App

1. Create a Notion page and place your own Notion DB.

- You need to create 5 columns in the Notion DB: `App ID` (Type: Title), `Title` (Type: Text), `Current Price` (Type: Number), `Lowest Price` (Type: Number), `Release Date` (Type: Date).

  ![Screenshot 2024-12-14 134649](https://github.com/user-attachments/assets/b9d65a3e-f15f-4d15-85c0-fa0194e96850)

2. Create an integration to use Notion API and connect it to the page where the Notion DB is set up.

- For Capabilities in the integration, you need to tick `Read content`, `Update content`, and `Insert content`.

3. Create your own Discord server and a Webhook.

4. Create a `.env` file.

   ```bash
    NOTION_API_KEY="dummy_notion_api_key"
    NOTION_DATABASE_ID="dummy_notion_database_id"
    DISCORD_WEBHOOK_ID="dummy_discord_webhook_id"
    DISCORD_WEBHOOK_TOKEN="dummy_discord_webhook_token"
    STEAM_USER_ID="dummy_steam_user_id"
   ```

5. Set up AWS infrastructure with AWS CDK.

   ```bash
    ./build.sh
    cd ./awscdk
    cdk deploy
   ```

6. Run the Lambda function manually by pressing Test

- When pressing the Test button on AWS Management Console, the app retrieves your Steam wishlist and write its data to the Notion DB.

7. Fill out the lowest prices of each game in the Notion DB.

- Please refer to [Steam DB](https://steamdb.info/) when filling out the lowest prices of each game in the Notion DB.

> [!IMPORTANT]
> Game prices are not notified unless you fill out the lowest prices of the video games in the Notion DB.
> Therefore, please do not skip Step 7 above.
