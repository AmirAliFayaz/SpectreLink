package main

import "SpectreLink/admin"

func main() {
	spec := admin.NewSpectreLink()
	spec.ListenAndServe()

}
