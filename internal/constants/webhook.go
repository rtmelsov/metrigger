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
var GaugeCommand = `
		INSERT INTO metrics (metric_name, metric_type, metric_value)
        VALUES ($1, $2, $3)
        ON CONFLICT (metric_name, metric_type) 
        DO UPDATE SET metric_value = EXCLUDED.metric_value;
	`
var CounterCommand = `
			INSERT INTO metrics (metric_name, metric_type, metric_delta)
        	VALUES ($1, $2, $3)
        	ON CONFLICT (metric_name, metric_type) 
        	DO UPDATE SET metric_delta = metrics.metric_delta + EXCLUDED.metric_delta;
		`
var GetGaugeRowCommand = `SELECT metric_value FROM metrics WHERE metric_type = $1 AND metric_name = $2`

var GetCounterRowCommand = `SELECT metric_delta FROM metrics WHERE metric_type = $1 AND metric_name = $2`

var GetRowCommand = `SELECT metric_value, metric_delta FROM metrics WHERE metric_type = $1 AND metric_name = $2`
