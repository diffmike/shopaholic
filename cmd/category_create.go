package cmd

import (
	"log"
	"shopaholic/store"
	"strings"
)

// CategoryCreateCommand set of flags and command for creation
type CategoryCreateCommand struct {
	Title string `short:"t" long:"title" description:"new category title" required:"true"`
	User  string `short:"u" long:"user" description:"ID of the user" required:"false"`

	CommonOpts
}

func (ccc *CategoryCreateCommand) Execute(args []string) error {
	log.Printf("[INFO] category %s creating command is started", ccc.Title)

	catType, catTitle := ccc.Store.DefineType(ccc.Title)
	category := store.Category{
		Title:  strings.Title(catTitle),
		UserID: ccc.User,
		Type:   catType,
	}

	categoryID, err := ccc.Store.StoreCategory(category)
	if err != nil {
		return err
	}

	log.Printf("[INFO] category %s was created with ID: %s", category.Title, categoryID)
	return nil
}
