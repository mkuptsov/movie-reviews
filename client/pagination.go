package client

import "github.com/mkuptsov/movie-reviews/contracts"

func Paginate[I any, Req contracts.PaginationSetter](
	req Req,
	requestFn func(Req) (*contracts.PaginatedResponse[I], error),
) ([]*I, error) {
	var items []*I

	for {
		res, err := requestFn(req)
		if err != nil {
			return nil, err
		}

		items = append(items, res.Items...)

		if len(items) >= res.Total {
			break
		}

		req.SetPage(res.Page + 1)
	}
	return items, nil
}
