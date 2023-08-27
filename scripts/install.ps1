$folders = @("golang", "python")

foreach ($folder in $folders) {
  if (Test-Path $folder) {
    Set-Location $folder

    switch ($folder) {
      "golang" {
        go mod tidy
      }
      "python" {
        pip install -r requirements.txt
      }
    }

    Set-Location .. 
  }
}