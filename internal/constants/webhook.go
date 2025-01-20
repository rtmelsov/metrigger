package constants

var Tmpl = `
<!DOCTYPE html>
			<html>
			<head>
				<title>Metrics</title>
			</head>
			<body>
				<div>
					<h3>Counter</h3>
					{{range $category, $product := .Counter}}
						<div>
							<h5>{{$category}}</h5>
							<ul>{{$product.Value}}</ul>
						</div>
					{{end}}
				</div>
				<div>
					<h3>Gauge</h3>
					{{range $category, $product := .Gauge}}
						<div>
							<h5>{{$category}}</h5>
							<ul>{{$product.Value}}</ul>
						</div>
					{{end}}
				</div>

			</body>
	</html>
	`
