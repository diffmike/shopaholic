package tg_bot

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"shopaholic/store"
	"strconv"
	"strings"
)

type CategoryCommand struct {
	Options
}

func (c *CategoryCommand) Execute(m *tb.Message) error {
	log.Printf("[INFO] running new category command with payload %s", m.Payload)

	catType, catTitle := c.Store.DefineType(m.Payload)
	category := store.Category{
		Title:  strings.Title(catTitle),
		UserID: strconv.Itoa(m.Sender.ID),
		Type:   catType,
	}

	categoryID, err := c.Store.StoreCategory(category)
	if err != nil {
		return err
	}

	log.Printf("[INFO] category %s was created with ID: %s", category.Title, categoryID)

	result := fmt.Sprintf("Category %s was created for you", category.Title)
	_, err = c.Bot.Send(m.Sender, result)

	return err
}
