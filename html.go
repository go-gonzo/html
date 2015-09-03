package html

// github.com/yosssi/gcss binding for gonzo.
// No Configuration required.

import (
	"bytes"
	"io/ioutil"

	"github.com/omeid/gonzo"
	"github.com/omeid/gonzo/context"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/svg"
)

type Options html.Options

func Minify(opt Options) gonzo.Stage {

	return func(ctx context.Context, in <-chan gonzo.File, out chan<- gonzo.File) error {

		for {
			select {
			case file, ok := <-in:
				if !ok {
					return nil
				}

				buff := new(bytes.Buffer)
				ctx.Infof("Minfiying %s", file.FileInfo().Name())
				m := minify.New()
				m.AddFunc("text/html", html.MinifyWithOptions(html.Options(opt)))
				m.AddFunc("text/css", css.Minify)
				m.AddFunc("text/javascript", js.Minify)
				m.AddFunc("image/svg+xml", svg.Minify)

				err := m.Minify("text/html", buff, file)
				if err != nil {
					return err
				}

				file = gonzo.NewFile(ioutil.NopCloser(buff), file.FileInfo())
				file.FileInfo().SetSize(int64(buff.Len()))

				out <- file
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
