# Capital Bank

This is a comman-line application for exchange information between ERP Capital and banks

## Sources
### PrivatBank
* [Опис API для взаємодії з серверною частиною Автоклієнта версія 3.0.0](https://docs.google.com/document/d/e/2PACX-1vTtKvGa3P4E-lDqLg3bHRF6Wi9S7GIjSMFEFxII5qQZBGxuTXs25hQNiUU1hMZQhOyx6BNvIZ1bVKSr/pub)
* [Інструкція “Автоклієнт”](https://docs.google.com/document/d/e/2PACX-1vS8rx2WKg69o6JvG5L4AhSXcU6vxXcJph6WK84qJcAYDBvsNYEob57jDMQhbosjc9gRS5bOTqTXf0vb/pub#h.nqpje6ikfhcq)

## Install
Run commang `go build`and copy to files to correct place for run. If you use Windows, you can create subfolder CapitalBank in "Program files" and put it there. For periodically running this application you can use SQL Agent with job step kind CmdExec and command line: `"c:\Program Files\CapitalBank\" && "capitalbank.exe"`

## Config
For correct work you must exactly define connection string to database. Application don't create a database structure, because was created to work with Capital database only.
