package docqueue

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/unullmass/msg-store/models"
)

var InsertChan = make(chan models.Document)
var QuitChan = make(chan os.Signal)

func InitWriter(db *gorm.DB) {

	retryChan := make(chan models.Document)

	for {
		select {
		case d := <-InsertChan:
			if err := db.Create(&d).Error; err != nil {
				// requeue for insertion
				retryChan <- d
			}
		case d := <-retryChan:
			if err := db.Create(d).Error; err != nil {
				log.Println("Error inserting Document ", d.ID)
			}
		case <-QuitChan:
			// flush all data
			for d := range InsertChan {
				if err := db.Create(&d).Error; err != nil {
					log.Println("Error inserting Document ", d.ID)
				}
			}
			// close the insert channel
			close(InsertChan)
			close(QuitChan)
			return
		}
	}
}
