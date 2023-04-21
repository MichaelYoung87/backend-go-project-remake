### `Installation/Information`
Om ni vill kunna testa programmet så måste ni använda er av:
https://console.cloud.google.com/ - Skapa konto och skapa projekt, buckets, SQL databas samt behörigheter för åtkomst.
Google Cloud Storage - Bucket - För att ladda upp profilbilder till.
Google Cloud SQL - Er databas för att spara information om användare samt profilbilds URL.

Eller ändra om själva i koden för att spara allting lokalt istället för att använda cloud tjänster.

Programmeringsspråket är Go - https://go.dev/dl/

Kodat i Visual Studio Code - https://code.visualstudio.com/download

## `Ni behöver både backend och frontend för att kunna köra mini-projektet.`

Frontend projektet finns på: https://github.com/MichaelYoung87/frontend-react-project-remake

Detta är instruktionerna för backend.

1. Skapa en mapp vart ni vill på datorn som heter t ex.  "backend-go-project-remake"
2. Navigera med terminal till denna mappen ni precis skapat t ex. om ni skapat mappen direkt på C: så skriver ni 'cd C:\backend-go-project-remake' utan ''.
3. När ni är i mappen 'C:\backend-go-project-remake>' Så skriver ni 'git clone https://github.com/MichaelYoung87/backend-go-project-remake' utan ''.
4. När det är klart så skriver ni 'code .' utan '' för att öppna VS Code med det hämtade projektet.
5. Kör även en 'go mod tidy' utan '' i terminalen för att ladda hem alla packages.
6. De .go filer och .JSON filer som behövs ändras är:
    authController.go - Rad 159 - Skriv in namnet på den bucket ni skapat i Google Cloud Storage
    connection.go - Rad 9 - Här måste ni skriva in username:password till er databas på Google Cloud SQL samt IP nummer och Portnummer (@tcp(XX.XXX.XX.XXX:3306)) och namnet på databasen (NAMN_PÅ_DATABAS) som ni skapat i Google Cloud SQL.

    NAMN_PÅ_JSON_HÄMTAD_FRÅN_GOOGLE_CLOUD.json - Denna hämtas från https://console.cloud.google.com/ som en .JSON fil. När ni skapat projekt, buckets, SQL databas och satt upp behörigheter för åtkomst så går ni in på projektet - trycker på hamburgare menyn uppe till vänster - går in på IAM & Admin - Service Accounts - Tryck på erat skapade projekt under 'Email' - Tryck på Keys fliken ovanför - Tryck på Add Key - Create new key - Välj att ni vill att key ska skapas som .JSON - Ladda ner denna .JSON fil och lägg sedan in den inuti keys\google\ mappen ni klonade projektet till t ex. backend-go-project-remake\keys\google\ - Den skall alltså ligga i samma mapp som NAMN_PÅ_JSON_HÄMTAD_FRÅN_GOOGLE_CLOUD.json ligger i nu.

    main.go - Rad 13 - os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "keys/google/NAMN_PÅ_JSON_HÄMTAD_FRÅN_GOOGLE_CLOUD.json") Här skriver ni då in namnet på den .JSON fil som ni lagt in i keys/google/ mappen.

7. i Go kör man från main.go, så när allt är korrekt installerat trycker ni på main.go och sedan trycker på "Play" ikonen uppe till höger Run Code (CTRL+ALT+N) för att starta igång backend. Backend måste vara igång för att frontenden ska fungera.
