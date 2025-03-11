# README

Models are plain struct-based types.
They are used only to define the data read from and written to the database.
This means that they can have sql.NullString, etc.
The consumer is expected to deal with it.

Use Cody to write the boilerplate code that maps query results to the Model structures.
