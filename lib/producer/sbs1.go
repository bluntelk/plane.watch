package producer

import (
	"bufio"
	"plane.watch/lib/tracker/sbs1"
)

func (p *producer) sbsScanner(scan *bufio.Scanner) error {
	for scan.Scan() {
		line := scan.Text()
		p.addFrame(sbs1.NewFrame(scan.Text()), &p.FrameSource)
		p.addDebug("SBS Frame: %s", line)
	}
	return scan.Err()
}
