## Use

```go
//	"go.mongodb.org/mongo-driver/mongo/options"
// "github.com/lixvyang/go-utils/bigdecimal"
opt := options.Client().
		SetRegistry(
			bson.NewRegistryBuilder().RegisterCodec(
				reflect.TypeOf(decimal.Decimal{}), &bigdecimal.MongoDecimal{},
			).Build()).
		SetDirect(true)
```