package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	ID   int    // Node ID
	Name string // Node Name
	X    int    // X coordinate
	Y    int    // Y coordinate
}

type Edge struct {
	Start int // Starting node ID of the edge
	End   int // Ending node ID of the edge
}

type Graph struct {
	Nodes       []Node        // List of nodes in the graph
	Edges       []Edge        // List of edges in the graph
	StartNodeID int           // ID of the start node
	EndNodeID   int           // ID of the end node
	AdjList     map[int][]int // Adjacency list representing the graph
}

// Function to print all nodes
func printNodes(nodes []Node) {
	fmt.Println("\nthe_rooms:")
	for _, node := range nodes {
		fmt.Printf("%d: %s (%d, %d)\n", node.ID, node.Name, node.X, node.Y)
	}
}

// Function to print all edges
func printEdges(edges []Edge) {
	fmt.Println("\nthe_links:")
	for _, edge := range edges {
		fmt.Printf("%d - %d\n", edge.Start, edge.End)
	}
}

// startNodeID'den endNodeID'ye kadar olan tüm yolları BFS kullanarak bulmak için bir fonksiyon
func (g *Graph) BFSAllPaths(startNodeID int, endNodeID int) [][]int {
	// Bulunan yolları saklamak için bir slice
	paths := [][]int{}
	// BFS kuyruğu, başlangıçta sadece başlangıç düğümünü içerir
	queue := [][]int{{startNodeID}}

	// Kuyruk boşalana kadar devam et
	for len(queue) > 0 {
		// Kuyruğun ilk yolunu al ve kuyruktan çıkar
		path := queue[0]
		queue = queue[1:]
		// Yolun son düğümünü al
		node := path[len(path)-1]

		// Eğer son düğüm bitiş düğümü ise, bu yolu sonuçlara ekle
		if node == endNodeID {
			paths = append(paths, path)
			continue
		}

		// Son düğümün komşularını kontrol et
		for _, neighbor := range g.AdjList[node] {
			// Eğer komşu zaten bu yolun içinde değilse, yeni bir yol oluştur
			if !contains(path, neighbor) {
				// Mevcut yolu kopyala
				newPath := append([]int{}, path...)
				// Yeni yolu komşu ile genişlet
				newPath = append(newPath, neighbor)
				// Yeni yolu kuyruğa ekle
				queue = append(queue, newPath)
			}
		}
	}

	// Bulunan tüm yolları geri döndür
	return paths
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Function to find the node ID by its name
// Bir düğümün adını kullanarak düğüm ID'sini bulmak için bir fonksiyon
func findNodeIDByName(nodes []Node, name string) int {
	// nodes slice'ındaki her bir node üzerinde döngü başlat
	for _, node := range nodes {
		// Eğer node'un adı verilen name'e eşitse
		if node.Name == name {
			// node.ID'yi döndür
			return node.ID
		}
	}
	// Eğer verilen name'e sahip bir node bulunamazsa, -1 döndür
	return -1
}

// Function to find an alternative path for an ant if the primary path is blocked
// Başlangıç düğümünden bitiş düğümüne giden birincil yol engellenmişse karınca için alternatif bir yol bulan bir fonksiyon.
func findAlternativePath(graph Graph, currentPos int, occupied map[int]bool) []int {
	// Tüm yolları hesaplamak için BFS kullanarak başlangıç konumundan bitiş düğümüne giden tüm yolları bul.
	allPaths := graph.BFSAllPaths(currentPos, graph.EndNodeID)

	// Tüm bulunan yollar üzerinde döngü başlat.
	for _, path := range allPaths {
		// Yolun geçerli olup olmadığını kontrol etmek için bir bayrak oluştur.
		valid := true

		// Yoldaki her düğüm için kontrol yap.
		for _, node := range path {
			// Eğer düğüm işgal edilmişse,
			if occupied[node] {
				// Yolu geçersiz olarak işaretle ve döngüyü sonlandır.
				valid = false
				break
			}
		}

		// Eğer yol geçerliyse,
		if valid {
			// Bu yolun alternatif bir yol olduğunu belirtmek için yol dizisini döndür.
			return path
		}
	}

	// Eğer engelsiz bir alternatif yol bulunamazsa, nil döndür.
	return nil
}

// FilterPaths, verilen yollar arasından düğüm çakışmalarını önleyerek en fazla sayıda yolu seçer.
func FilterPaths(paths [][]int) [][]int {
	// maxPaths, en fazla sayıda geçerli yolu saklar.
	var maxPaths [][]int
	// currentPaths, geçerli durumda incelenen yolları saklar.
	var currentPaths [][]int
	// usedNodes, kullanılan düğümleri izlemek için bir harita.
	usedNodes := make(map[int]bool)

	// backtrack, geriye izleme algoritması için iç içe bir fonksiyon olarak tanımlanır.
	var backtrack func(int)
	backtrack = func(start int) {
		// Eğer geçerli yolların sayısı, maksimum yollardan fazlaysa, maxPaths güncellenir.
		if len(currentPaths) > len(maxPaths) {
			maxPaths = make([][]int, len(currentPaths)) //maxPaths değişkeni, currentPaths ile aynı sayıda elemana sahip boş bir slice'e dönüşüyor
			//maxpaths değişkeni adında currentpaths uzunluğunda int değerinde değişken oluşturuyor
			copy(maxPaths, currentPaths) //içerisine kopyalıyor
		}

		// Başlangıç indeksinden yolların sonuna kadar dolaş.
		for i := start; i < len(paths); i++ {
			path := paths[i] //yolları tek tek path değişkenine atıyor
			keepPath := true //kullanılma durumunu kontrol ediyor

			// İlk ve son düğüm hariç, yolun düğümlerini kontrol et.
			for _, node := range path[1 : len(path)-1] { //İlk ve son değişkenleri hariç bütün değişkenler alsın
				// Eğer düğüm daha önce kullanıldıysa, bu yolu kullanma.
				if usedNodes[node] {
					keepPath = false
					break
				}
			}

			// Eğer yol geçerliyse (düğümler kullanılmamışsa), yollar listesine ekle.
			if keepPath {
				currentPaths = append(currentPaths, path)
				// Kullanılan düğümleri işaretle.
				for _, node := range path[1 : len(path)-1] { //baştaki ve sondaki eleman hariç elemanlar üstünde dolaş
					usedNodes[node] = true //daha önce kullanıldı olarak değişiklik yapar
				}

				// Bir sonraki yol için geriye izleme (backtracking) yap.
				backtrack(i + 1)
				/*Bu işlem, bir yolun tamamlanmasından sonra diğer olası
				yolları aramak için tekrarlanır, böylece tüm olası yollar taranır ve en uzun, üst üste binmeyen yollar bulunur.*/
				//yeni bir yol arayışını başlatır
				// Backtrack: Son eklenen yolu ve düğümleri geri al.
				currentPaths = currentPaths[:len(currentPaths)-1] //Bu adımlar, geri izleme işlemi sırasında, bir sonraki olası yolu aramak için bir önceki adıma geri dönülmesini sağlar.
				for _, node := range path[1 : len(path)-1] {      //geri dönerek ihtimallerini buluyor
					delete(usedNodes, node)
				}
			}
		}
	}

	// Geriye izleme algoritmasını başlat.
	backtrack(0)
	return maxPaths
}

func main() {
	startTime := time.Now() // Başlangıç zamanını al

	if len(os.Args) != 2 {
		fmt.Println("Dosya adı belirtilmedi.")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Dosya açma hatası:", err)
		return
	}
	defer file.Close()

	graph := Graph{
		AdjList:     make(map[int][]int), // Düğümlerin komşuluk ilişkilerini depolamak için bir harita oluşturulur.
		StartNodeID: -1, // Başlangıç düğümünün ID'si -1 olarak başlatılır (bu değer daha sonra belirlenecektir).
		EndNodeID:   -1, // Bitiş düğümünün ID'si -1 olarak başlatılır (bu değer daha sonra belirlenecektir).
	}
	

	scanner := bufio.NewScanner(file)

// Sayıda karıncayı oku
scanner.Scan()
antCountLine := scanner.Text()
antCount, err := strconv.Atoi(antCountLine)
if err != nil {
    fmt.Println("HATA: Geçersiz veri formatı") // Hata durumunda geçersiz veri formatı hatası yazdırılır
    return
}
if antCount <= 0 {
    fmt.Println("HATA: Geçersiz veri formatı") // Hata durumunda geçersiz veri formatı hatası yazdırılır
    return
}


	// Graf verilerini oku
for scanner.Scan() {
    line := scanner.Text() // Bir sonraki satırı oku
    if strings.HasPrefix(line, "##start") {
        scanner.Scan()
        fields := strings.Fields(scanner.Text())
        if len(fields) < 3 {
            fmt.Println("HATA: Geçersiz veri formatı")
            return
        }
        name := fields[0] // Oda ismini al
        x, err := strconv.Atoi(fields[1]) // X koordinatını al
        if err != nil {
            fmt.Println("HATA: Geçersiz veri formatı")
            return
        }
        y, err := strconv.Atoi(fields[2]) // Y koordinatını al
        if err != nil {
            fmt.Println("HATA: Geçersiz veri formatı")
            return
        }
        startID := len(graph.Nodes) // Başlangıç düğümünün ID'sini belirle
        graph.StartNodeID = startID // Graf yapısındaki başlangıç düğüm ID'sini güncelle
        graph.Nodes = append(graph.Nodes, Node{ID: startID, Name: name, X: x, Y: y}) // Başlangıç düğümünü graf düğümlerine ekle

    } else if strings.HasPrefix(line, "##end") {
        scanner.Scan()
        fields := strings.Fields(scanner.Text())
        if len(fields) < 3 {
            fmt.Println("HATA: Geçersiz veri formatı")
            return
        }
        name := fields[0] // Oda ismini al
        x, err := strconv.Atoi(fields[1]) // X koordinatını al
        if err != nil {
            fmt.Println("HATA: Geçersiz veri formatı")
            return
        }
        y, err := strconv.Atoi(fields[2]) // Y koordinatını al
        if err != nil {
            fmt.Println("HATA: Geçersiz veri formatı")
            return
        }
        endID := len(graph.Nodes) // Bitiş düğümünün ID'sini belirle
        graph.EndNodeID = endID // Graf yapısındaki bitiş düğüm ID'sini güncelle
        graph.Nodes = append(graph.Nodes, Node{ID: endID, Name: name, X: x, Y: y}) // Bitiş düğümünü graf düğümlerine ekle

    } else {
        fields := strings.Fields(line) // Satırı alanlara ayır
        if len(fields) == 3 { // Eğer üç alana ayrılmışsa (ID, X, Y)
            name := fields[0] // Oda ismini al
            x, err := strconv.Atoi(fields[1]) // X koordinatını al
            if err != nil {
                fmt.Println("HATA: Geçersiz veri formatı")
                return
            }
            y, err := strconv.Atoi(fields[2]) // Y koordinatını al
            if err != nil {
                fmt.Println("HATA: Geçersiz veri formatı")
                return
            }
            id := len(graph.Nodes) // Düğümün ID'sini belirle
            graph.Nodes = append(graph.Nodes, Node{ID: id, Name: name, X: x, Y: y}) // Düğümü graf düğümlerine ekle
        } else if len(fields) == 1 && strings.Contains(line, "-") { // Eğer bir alan içeriyor ve içinde "-" karakteri varsa (bir kenar)
            edgeParts := strings.Split(fields[0], "-") // Kenarı ayır
            if len(edgeParts) != 2 { // Eğer iki kısım yoksa (başlangıç ve bitiş düğümleri eksikse)
                fmt.Println("HATA: Geçersiz veri formatı")
                return
            }
            startName := edgeParts[0] // Başlangıç düğüm ismini al
            endName := edgeParts[1] // Bitiş düğüm ismini al
            startID := findNodeIDByName(graph.Nodes, startName) // Başlangıç düğüm ID'sini bul
            endID := findNodeIDByName(graph.Nodes, endName) // Bitiş düğüm ID'sini bul
            if startID == -1 || endID == -1 { // Eğer başlangıç veya bitiş düğümü bulunamadıysa
                fmt.Println("HATA: Geçersiz veri formatı")
                return
            }
            graph.Edges = append(graph.Edges, Edge{Start: startID, End: endID}) // Kenarı graf kenarlarına ekle
            graph.AdjList[startID] = append(graph.AdjList[startID], endID) // Başlangıç düğümünün komşuları listesine bitiş düğümünü ekle
            graph.AdjList[endID] = append(graph.AdjList[endID], startID) // Bitiş düğümünün komşuları listesine başlangıç düğümünü ekle
        }
    }
}


if err := scanner.Err(); err != nil {
	fmt.Println("HATA: Okuma hatası") // Okuma hatası
	return
}

if graph.StartNodeID == -1 || graph.EndNodeID == -1 {
	fmt.Println("HATA: Başlangıç veya bitiş düğümü belirtilmedi") // Hata: Başlangıç veya bitiş düğümü belirtilmedi
	return
}

allPaths := graph.BFSAllPaths(graph.StartNodeID, graph.EndNodeID)
if len(allPaths) == 0 {
	fmt.Println("HATA: Geçersiz veri formatı")
	return
}

// Tüm yolları, uzunluklarına göre sıralar.
sort.Slice(allPaths, func(i, j int) bool {
	return len(allPaths[i]) < len(allPaths[j])
})

// Sıralanmış yollardan oluşan 'allPaths' slice'ını filtreler.
filteredPaths := FilterPaths(allPaths)

// Eğer filtrelenmiş yolların uzunluğu 0 ise veya ilk yolu boşsa, 
// bu durumda geçersiz bir veri formatı olduğunu belirten bir hata mesajı yazdırılır.
if len(filteredPaths) == 0 || len(filteredPaths[0]) == 0 {
	fmt.Println("HATA: Geçersiz veri formatı")
	return
}


// Giriş verilerini yazdır
fmt.Printf("Karınca sayısı: %d\n", antCount)
fmt.Printf("Başlangıç odası: %d\n", graph.StartNodeID)
fmt.Printf("Bitiş odası: %d\n", graph.EndNodeID)
printNodes(graph.Nodes)
printEdges(graph.Edges)

antPaths := assignPathsToAnts(antCount, filteredPaths)
antPositions := make([]int, antCount)
antAtEnd := make([]bool, antCount)

// Bütün karıncaların pozisyonlarını ve yollarını başlat
for i := 0; i < antCount; i++ {
	if i == antCount-1 {
		antPaths[i] = filteredPaths[0] // Son karınca en kısa yolu takip eder
	} else {
		antPaths[i] = filteredPaths[(i)%len(filteredPaths)] // Diğer karıncalar sırayla takip eder
	}
	antPositions[i] = graph.StartNodeID // Bütün karıncalar başlangıç pozisyonundadır
	antAtEnd[i] = false                 // Hiçbir karınca bitişe ulaşmamıştır
}


step := 1 // Adım sayacı başlatılır.

// Sonsuz döngü başlatılır. Döngü, tüm karıncaların hedefe ulaşıncaya kadar devam eder.
for {
	// Karıncaların yapacağı hareketlerin listesi başlatılır.
	moves := []string{}

	// Tüm karıncaların hedefe ulaşıp ulaşmadığı kontrol edilir.
	allAtEnd := true

	// Karıncaların mevcut konumları takip edilir.
	occupied := make(map[int]bool)

	// Bu adım için planlanmış hareketler toplanır.
	for i := 0; i < antCount; i++ {
		// Eğer karınca hedefe ulaşmadıysa, döngü devam eder.
		if antPositions[i] != graph.EndNodeID {
			allAtEnd = false // En az bir karıncanın hedefe ulaşmadığını belirtmek için bayrak ayarlanır.
			
			// Karıncanın takip ettiği yol alınır.
			path := antPaths[i]
			for j := 0; j < len(path)-1; j++ {
				// Karıncanın mevcut konumu ile hedefi arasında bir bağlantı var mı kontrol edilir.
				if path[j] == antPositions[i] && (!occupied[path[j+1]] || path[j+1] == graph.EndNodeID) {
					// Eğer bir bağlantı varsa, karıncanın yeni konumu güncellenir ve bu hareket kaydedilir.
					occupied[path[j+1]] = true
					antPositions[i] = path[j+1]
					moves = append(moves, fmt.Sprintf("L%d-%s", i+1, graph.Nodes[antPositions[i]].Name))
					break
				}
			}

			// Eğer karıncanın yolu engellenmişse, alternatif yol aranır.
			if antPositions[i] != graph.EndNodeID && occupied[antPositions[i]] {
				altPath := findAlternativePath(graph, antPositions[i], occupied)
				if altPath != nil {
					antPaths[i] = altPath // Alternatif yol bulunursa, karıncanın yolu güncellenir.
				}
			}
		} else {
			antAtEnd[i] = true // Karınca hedefe ulaştıysa, bu bilgi kaydedilir.
		}
	}

	// Bu adımda yapılacak hareketler ekrana yazdırılır.
	if len(moves) > 0 {
		fmt.Printf("Adım %d: %s\n", step, strings.Join(moves, " "))
	}

	// Eğer tüm karıncalar hedefe ulaştıysa, döngüden çıkılır.
	if allAtEnd {
		break
	}

	step++ // Adım sayacı artırılır.

	// Eğer tüm karıncaların son adımda hedefe ulaşamadığı tespit edilirse, döngüden çıkılır.
	allAtEnd = true
	for i := 0; i < antCount; i++ {
		if !antAtEnd[i] {
			allAtEnd = false
			break
		}
	}
}
// elapsed time hesaplanır.
elapsedTime := time.Since(startTime)

// Toplam geçen süre saniye cinsinden hesaplanır ve kesirli kısmı ile birlikte ekrana yazdırılır.
fmt.Printf("Toplam süre: %.9f saniye\n", elapsedTime.Seconds())
}

// Karıncalara yol atamak için kullanılan fonksiyon. Yolların örtüşmemesini sağlar ve ilk karınca için en kısa yolu seçer.
func assignPathsToAnts(antCount int, filteredPaths [][]int) [][]int {
	antPaths := make([][]int, antCount)  // Karıncaların yollarını depolamak için bir slice oluşturulur.
	remainingPaths := filteredPaths      // Atanmamış yolları depolamak için bir değişken oluşturulur.

	// Tüm karıncalara yollar atanır
	for i := 0; i < antCount; i++ {
		if len(remainingPaths) > 0 {
			antPaths[i] = remainingPaths[0]  // Atanmamış yollardan biri bu karıncaya atanır.
			remainingPaths = remainingPaths[1:]  // Atanmış olan yolu listeden kaldırılır.
		} else {
			// Eğer daha fazla benzersiz yol yoksa, mevcut yollar tekrar kullanılır.
			antPaths[i] = filteredPaths[i%len(filteredPaths)]
		}
	}

	return antPaths  // Tüm karıncaların atandığı yolların bulunduğu slice döndürülür.
}