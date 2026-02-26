package mongodb

import (
	"dungeons/app/models"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// SelectConstructeur : Builds the select request
func SelectConstructeur(params models.QueryParams) bson.M {

	var findQuery bson.M
	findQuery = make(map[string]interface{})
	if params.TestDeleted {
		findQuery = Alive(findQuery)
	}
	findQuery = FilterConstructeur(params, findQuery)
	findQuery = FilterLikeConstructeur(params, findQuery)
	return findQuery

}

// FilterConstructeur : Preparation du where
func FilterConstructeur(params models.QueryParams, fq bson.M) bson.M {
	var (
		key      string
		value    interface{}
		operator string
		query    bson.M
	)

	if len(params.FilterClause) > 0 {
		for _, param := range params.FilterClause {

			// Retrieve the requested filter
			filter := strings.Split(param, ",")

			// Operator retrieve
			if len(filter) == 2 {
				operator = "="
			} else {
				operator = filter[2]
			}

			key = filter[0]

			// conversion of the value into its data type
			if vBool, sErr := strconv.ParseBool(filter[1]); sErr != nil {
				if vInt, sErr := strconv.ParseInt(filter[1], 10, 64); sErr != nil {
					if vFloat, sErr := strconv.ParseFloat(filter[1], 64); sErr != nil {
						value = filter[1]
					} else {
						value = vFloat
					}
				} else {
					value = vInt
				}
			} else {
				value = vBool
			}

			// operator management
			switch operator {
			case "=":
				fq[key] = value
			case ">", ">=", "<", "<=", "!=":
				if query == nil {
					query = make(bson.M) // Initialise query comme une map vide
				}

				switch operator {
				case ">":
					query["$gt"] = value
				case ">=":
					query["$gte"] = value
				case "<":
					query["$lt"] = value
				case "<=":
					query["$lte"] = value
				case "!=":
					query["$ne"] = value
				}
				fq[key] = query
			default:
				fq[key] = value
			}

		}
	}
	return fq
}

// FilterLikeConstructeur : preparing where for a like
func FilterLikeConstructeur(params models.QueryParams, fq bson.M) bson.M {

	if len(params.FilterLikeClause) > 0 {
		for _, param := range params.FilterLikeClause {
			// Retrieving the requested like filter
			filterLike := strings.Split(param, ",")
			fq[filterLike[0]] = primitive.Regex{Pattern: "^.*" + filterLike[1] + ".*$", Options: "i"}
		}
	}
	return fq
}

// Alive allows to set the element to deleted (Archive)
func Alive(fq bson.M) bson.M {
	fq["deleted"] = false
	return fq
}
