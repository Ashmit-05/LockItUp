package connectToDatabase

import (
  "log"
  "os"
  "github.com/joho/godotenv"
)
func connectToDB()  {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error in reading .env file")
  }

  connectionString := os.Getenv("MONGODB_URI")
  
}
