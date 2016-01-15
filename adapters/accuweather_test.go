package adapters

import "testing"

func TestAccuweatherIsMetric(t *testing.T) {
	s := `<div id="control-panel"><div id="social-connect"><ul class="social-list"><li class="label">Следите за нашими сообщениями в</li><li class="social-icon"><a id="connect_googleplus" href="https://plus.google.com/+accuweather/posts" class="gplus"></a></li><li class="social-icon"><a id="connect_youtube" href="https://www.youtube.com/user/accuweather" class="youtube"></a></li><li class="social-icon"><a id="connect_twitter" href="https://twitter.com/breakingweather" class="twitter"></a></li><li class="social-icon"><a id="connect_facebook" href="https://www.facebook.com/AccuWeather" class="facebook"></a></li></ul></div><a id="bt-menu-login" class="tmenu { el:'#menu-premium', affix: { to: 'ne', from: 'ne', offset: [ 0,-6 ] } }"><span class="menu-arrow"><span>Логин</span></span></a> <a id="bt-menu-settings" class="tmenu { el:'#menu-settings', affix: { from: 'ne', to: 'ne', offset: [ 0,-6 ] } }"><span class="menu-arrow"><span>русский</span>, °C</span></a></div>`
	expectation := true

	result, resultErr := AccuweatherIsMetric(s)

	if resultErr != nil {
		t.Errorf(resultErr.Error())
	}

	if result != expectation {
		t.Errorf(ErrorOut(expectation, result))
	}
}
