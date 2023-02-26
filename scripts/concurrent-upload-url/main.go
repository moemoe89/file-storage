package main

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
)

func main() {
	images := []string{
		"https://images.ctfassets.net/30xxrlh9suih/6jyUVNJs0bKkYwVGAxarUa/1fa3440d4c578d6ce05d38dbbf9fcb08/wovenplanet_60px_onwhite.jpeg",
		"https://images.ctfassets.net/30xxrlh9suih/1swYJcv6FRXiMa4WemPjnI/b5836346de2f2839ce79be4eb1c2a527/Arene_Hex.jpg",
		"https://images.ctfassets.net/30xxrlh9suih/7ImuEjQRqe7KNxKbecVHjf/f5be9df280ba29fea232cd83fab1ada6/Teammate_HEX.jpg",
		"https://images.ctfassets.net/30xxrlh9suih/79goF34bhyagSdlRApqgUb/df33fee8a77629817336656b3575de4c/WP_P4_sideview.jpg",
		"https://images.ctfassets.net/30xxrlh9suih/3cLSbThmHt3znyW9lAkWLB/813977aee3cce73d6683753f7b020baf/automated-mapping.jpg",
		"https://images.ctfassets.net/30xxrlh9suih/7HfDVTVGFtPNEPTUr9QLnR/d48cdde2b5de54dba9e35295f2ace702/Woven_City.png",
		"https://www.woven-planet.global/assets/images/v2/car.png",
		"https://images.ctfassets.net/30xxrlh9suih/3jekS0uKNjqE8sUU5gkPmX/fe9356893e9e848bf14fb00201c987d8/Minh.png",
		"https://images.ctfassets.net/30xxrlh9suih/4A0igpF3ADzwI8EgXdmS3Y/95d986560dc567b08612050aefc40802/Chris.png",
		"https://images.ctfassets.net/30xxrlh9suih/148pyTqTd3wWyBFOgUj3si/1298046c36aac24f08e8982a448d185f/Meredith.png",
		"https://images.ctfassets.net/30xxrlh9suih/54dUKFZrnMb4g4cmGR6nMZ/68a3b9cdee906c34065e6e60fb1adb2b/Jack.png",
		"https://images.ctfassets.net/30xxrlh9suih/1XEZTwwPKeSnKwIhc2CrcS/8e8961e00eb1fbc8805c3d16edbbb0d4/employability_1.JPG",
		"https://images.ctfassets.net/30xxrlh9suih/34hxnlkB5U7frgst11y5kw/535bc114ec72c26535d9398c7611d057/employability_2.jpg",
		"https://images.ctfassets.net/30xxrlh9suih/1kvBOT1Jo4gSPQpM9Qi7WK/3a6b0cdf4b2e6b19e7f7b06b65e92a04/Tokyo.png",
		"https://images.ctfassets.net/30xxrlh9suih/1ihGB2Vg4JeeDmbHJdwayk/5c8894cb798b279d9b851988834aca83/Bay_Area.png",
		"https://images.ctfassets.net/30xxrlh9suih/19Gbkg83MRt0ObHZsD58Wf/0011a208c5051d8ce5c8c21a0742e218/London.png",
		"https://images.ctfassets.net/30xxrlh9suih/3zZztbggwX0ebhKku0JZZK/dabeb5d9af7cb6b63aa67afa7039e511/Seattle.png",
		"https://images.ctfassets.net/30xxrlh9suih/2FqNxPuC9Oh6vXuggtNFt1/05151878658aed5f04bcfae35ab5f74e/value_1.jpg",
		"https://images.ctfassets.net/30xxrlh9suih/6NdtchTHaZMLWuB9wboSj0/b07ddeb1192645c7e735796ffae64a99/value_2.jpg",
		"https://images.ctfassets.net/30xxrlh9suih/4cGcmdHK08U8O2f6Kl3dKc/434bcd718603eba0cf7733cff86d71d8/value_3.jpg",
		"https://images.ctfassets.net/30xxrlh9suih/oh1zwt55pBjUogj1cNwki/4d9816efe0ea0cebdad8fe640dc47b73/diversity-new.jpg",
	}

	// Non-concurrent test
	//for _, image := range images {
	//  cmd := "./fs-store -source=url upload-file " + image
	//
	//  _, err := exec.Command("bash", "-c", cmd).Output()
	//	if err != nil {
	//		log.Println(err.Error())
	//	}
	//
	//	fmt.Println(string(cmd))
	//}

	errChan := make(chan error, len(images))

	var wg sync.WaitGroup

	for _, image := range images {
		wg.Add(1)

		go func(image string) {
			defer wg.Done()

			cmd := "./fs-store -source=url upload-file " + image

			_, err := exec.Command("bash", "-c", cmd).Output()
			if err != nil {
				errChan <- fmt.Errorf("failed to upload: %w", err)
			}

		}(image)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		log.Fatal(err)
	default:
		close(errChan)
		log.Printf("### Uploading finished!! ###\n")
	}
}
