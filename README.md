# Capital Bank

This is a comman-line application for exchange information between ERP Capital and banks
## Features
* implement API requests to Privatbank for take transactions list and balances
* implement processing .csv files in iBank2UA specification

## Sources
### PrivatBank
* [Опис API для взаємодії з серверною частиною Автоклієнта версія 3.0.0](https://docs.google.com/document/d/e/2PACX-1vTtKvGa3P4E-lDqLg3bHRF6Wi9S7GIjSMFEFxII5qQZBGxuTXs25hQNiUU1hMZQhOyx6BNvIZ1bVKSr/pub)
* [Інструкція “Автоклієнт”](https://docs.google.com/document/d/e/2PACX-1vS8rx2WKg69o6JvG5L4AhSXcU6vxXcJph6WK84qJcAYDBvsNYEob57jDMQhbosjc9gRS5bOTqTXf0vb/pub#h.nqpje6ikfhcq)
### iBank2UA csv
folder docs/

## Install
Run command `go build`and copy to files to correct place for run. If you use Windows, you can create subfolder CapitalBank in "Program files" and put it there. For periodically running this application you can use SQL Agent with job step kind CmdExec and command line: `"c:\Program Files\CapitalBank\" && "capitalbank.exe"`

## Config
For correct work you must exactly define connection string to database. To work with sql instances you should define server name as "name\\instancename"
Application dosn't create a database structure, because was created to work with Capital database only.

## Program structure
api - interfaces, which should be defined for any new bank
logger - define logger functionality
logic - define a business logic program's behaviour
pbapi - methods defined for Privatbank processing
store - methods for store data in the database
utils - some additional general functions

## Algorithm

### main idea
* load list of bank accounts 
* take data from banks
* converting data to one format for saving
* save data in the database

### details
* main  cycle in folder logic/
* work with csv-files has trick: all files should be saved to one catalog; this files should be loaded before render and convert to one structure; after processing all files are deleted
* csv format doesn't have a statements data - transactions only
* any new bank format must implement all methods of BankAPI Interface 
* some banks doesn't have an unique ID of transaction, we need find unique combination before saving to DB

