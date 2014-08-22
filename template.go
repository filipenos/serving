package main

var (
	upload = `<!DOCTYPE html>
				<html lang="en">
				<head>
				<title>File upload</title>
				</head>
				<body>
				<div>
					<h1>File Upload</h1>
					<form method="post" action="/upload" enctype="multipart/form-data">
						<fieldset>
							<input type="file" name="files" id="files" multiple="multiple">
							<input type="submit" name="submit" value="Submit">
						</fieldset>
					</form>
				</div>
				</body>
				</html>`
)
