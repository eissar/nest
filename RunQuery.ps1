Param(
    [Parameter(Mandatory=$true)]
    $DatabaseConnection
    [String]$QueryFile,
    $Query
)

if ($QueryFile) {
    Assert-PathExists -FilePath $QueryFile
    cat $QueryFile | sqlite3.exe -cmd '.open test.sqlite'
} else {
    Write-Output $Query | sqlite3.exe

}


