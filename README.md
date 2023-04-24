# TD-EXPORT v0.1.0
MySQL Table Definition Export tool.

Export to Excel and Markdown is supported.

## Usang
```shell
./TD-EXPORT -output=[excel|markdown]
```

### Example
```shell
 % ./TD-EXPORT -output=markdown
INFO[0000] TD-EXPORT v0.1.0                             
INFO[0000] Table Definition Export                      

Endpoint : test1.ap-northeast-2.rds.amazonaws.com
Port : 
WARN[0006] Not Port Args.                               
INFO[0006] Default set Port : 3306                      
User : root
Password : 
DB(Seperator , or Space(All)) : db1,db2
Exception Tables(Seperator , or Space(none) / Use wildcard) : %_tale1,%_table2
INFO[0027] DB Connect Success                           
INFO[0027] Setup Markdown Files                         
INFO[0027] Get Schema Count : 2                        
INFO[0027] db1 Table Load.                       
INFO[0027] db1 Table Count : 10                 
INFO[0027] db1 Table Column/Index/Const Load   
INFO[0028] db2 Table Load.                       
INFO[0028] db2 Table Count : 10                 
INFO[0028] db2 Table Column/Index/Const Load   
INFO[0031] Export Complete. 
```
- When selecting multiple DBs, separate them with commas (,). Excel divides them into sheets, and Markdown divides them into files.
- Exception table division is separated by a comma (,). If you select multiple exception tables, extract for all tables except the target table.
- The extracted file is saved in the same path as the tool, and the file name is as follows.
    - Excel : ```{Endpoint}.xlsx```
    - Markdown : ```{Schema}({Endpoint}).md```