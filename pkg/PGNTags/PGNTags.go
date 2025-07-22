package pgntags

/*
	[Event "F/S Return Match"]
	[Site "Belgrade, Serbia JUG"]
	[Date "1992.11.04"]
	[Round "29"]
	[White "Fischer, Robert J."]
	[Black "Spassky, Boris V."]
	[Result "1/2-1/2"]
*/

type PGNTags struct {
	Event  string
	Site   string
	Date   string
	Round  int
	White  string
	Black  string
	Result string
}
