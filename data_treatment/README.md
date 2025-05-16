# Go data_treatment

I created the file `resources/config.conf` to do quicker configurations.

For using this code, you must add to `resources/mail.conf` the correspondent data (if not, it won't work):

```text
MAIL_ACCOUNT=example@gmail.com
MAIL_PASSWORD=example_password
SMTP_SERVER=smtp.gmail.com
```

Also, you must add to `resources/send_to.conf` the receivers:

```text
SEND_TO=example2@gmail.com
```

## Important

You must have enabled the 2-steps verification in the account.

Nowadays, for using this, we have to create an App pasword from Google account settings.

For Google: you enter in My Account of Google -> Security -> 2-Step Verification -> App Passwords.

You must create 1 and copy to the .env file (in the MAIL_PASSWORD field, with no spaces).

## To execute this code

```bash
chmod +x launch.sh
./launch.sh <connection_string>
```
