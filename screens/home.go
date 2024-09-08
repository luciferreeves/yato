package screens

import (
	"fmt"
	"yato/config"
	"yato/lib"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomeScreen struct {
	RecentAnimeRecommendations []lib.Recommendation
	RecentMangaRecommendations []lib.Recommendation
	imageCache                 *lib.ImageCache
	imageRenderer              *lib.ImageRenderer
}

func homeScreen() tea.Model {
	recentAnimeRecommendations, _ := lib.GetRecentAnimeRecommendations()
	recentMangaRecommendations, _ := lib.GetRecentMangaRecommendations()
	imageCache := lib.NewImageCache()
	imageRenderer := lib.NewImageRenderer()

	return HomeScreen{
		RecentAnimeRecommendations: recentAnimeRecommendations,
		RecentMangaRecommendations: recentMangaRecommendations,
		imageCache:                 imageCache,
		imageRenderer:              imageRenderer,
	}
}

func (h HomeScreen) Init() tea.Cmd {
	return nil
}

func (h HomeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.String() == "q" {
			return h, tea.Quit
		}
	}

	return h, nil
}

func (h HomeScreen) View() string {
	w := lipgloss.Width

	// Top bar, Content, Status bar
	topBarStyle := lipgloss.NewStyle().
		Foreground(config.Colors.Text).
		Background(config.Colors.Primary)

	mainText := topBarStyle.Padding(0, 0, 0, 1).Render(config.PrettyAppName + " | [H]ome | [A]nime | [M]anga | [S]earch | [C]ommunity | [P]rofile | [O]ptions | [Q]uit")
	userText := topBarStyle.Padding(0, 1, 0, 0).Render("User: " + globals.CurrentUser.Name + " | [L]ogout")
	separator := topBarStyle.Width(globals.width - w(mainText) - w(userText)).Render("")

	topBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		mainText,
		separator,
		userText,
	)

	content := ""
	// Top 5 recommendations
	for i, rec := range h.RecentAnimeRecommendations {
		if i == 5 {
			break
		}
		content += h.renderRecommendation("anime", rec)
	}

	for i, rec := range h.RecentMangaRecommendations {
		if i == 5 {
			break
		}
		content += h.renderRecommendation("manga", rec)
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		topBar,
		content,
	)

}

func (h HomeScreen) renderRecommendation(mediaType string, rec lib.Recommendation) string {
	img, err := h.imageCache.GetImage(mediaType, rec.Entry[0].MALId, "small", rec.Entry[0].Images.JPG.SmallImageURL)
	if err != nil {
		return fmt.Sprintf("%s -> %s\n", rec.Entry[0].Title, rec.Entry[1].Title)
	}

	renderedImage := h.imageRenderer.RenderImage(img, 20, 30)
	return fmt.Sprintf("%s%s -> %s\n", renderedImage, rec.Entry[0].Title, rec.Entry[1].Title)
}
