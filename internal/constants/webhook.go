package constants

var Tmpl = `
<!DOCTYPE html>
			<html>
			<head>
				<title>Metrics</title>
<style>
table {
  font-family: arial, sans-serif;
  border-collapse: collapse;
  width: 100%;
}

td, th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}

tr:nth-child(even) {
  background-color: #dddddd;
}
</style>
			</head>
			<body>
				<div>
					<h3>Counter</h3>


						<table>
					{{range $category, $product := .Counter}}
						  <tr>
							<th>{{$category}}</th>
							<th>{{$product.Value}}</th>
						  </tr>
					{{end}}

						</table>
				</div>
				<div>
					<h3>Gauge</h3>
						<table>
					{{range $category, $product := .Gauge}}
					      <tr>
							<th>{{$category}}</th>
							<th>{{$product.Value}}</th>
						  </tr>
					{{end}}

						</table>
				</div>

			</body>
	</html>
	`
