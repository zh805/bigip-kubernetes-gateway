package pkg

import (
	"sync"

	"gitee.com/zongzw/f5-bigip-rest/utils"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

func init() {
	PendingDeploys = make(chan DeployRequest, 16)
	PendingParses = make(chan ParseRequest, 16)
	slog = utils.SetupLog("", "debug")
	ActiveSIGs = &SIGCache{
		mutex:     sync.RWMutex{},
		Gateway:   map[string]*gatewayv1beta1.Gateway{},
		HTTPRoute: map[string]*gatewayv1beta1.HTTPRoute{},
	}
}

func (c *SIGCache) SetGateway(obj *gatewayv1beta1.Gateway) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if obj != nil {
		c.Gateway[utils.Keyname(obj.Namespace, obj.Name)] = obj
	}
}

func (c *SIGCache) UnsetGateway(keyname string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.Gateway, keyname)
}

func (c *SIGCache) GetGateway(keyname string) *gatewayv1beta1.Gateway {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.Gateway[keyname]
}

func (c *SIGCache) SetHTTPRoute(obj *gatewayv1beta1.HTTPRoute) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if obj != nil {
		c.HTTPRoute[utils.Keyname(obj.Namespace, obj.Name)] = obj
	}
}

func (c *SIGCache) UnsetHTTPRoute(keyname string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.HTTPRoute, keyname)
}

func (c *SIGCache) GetHTTPRoute(keyname string) *gatewayv1beta1.HTTPRoute {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.HTTPRoute[keyname]
}

func (c *SIGCache) ParentRefsOf(hr *gatewayv1beta1.HTTPRoute) []*gatewayv1beta1.Gateway {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if hr == nil {
		return []*gatewayv1beta1.Gateway{}
	}
	gws := []*gatewayv1beta1.Gateway{}
	for _, pr := range hr.Spec.ParentRefs {
		ns := hr.Namespace
		if pr.Namespace != nil {
			ns = string(*pr.Namespace)
		}
		name := utils.Keyname(utils.Keyname(ns, string(pr.Name)))
		if gw, ok := c.Gateway[name]; ok {
			gws = append(gws, gw)
		}
	}
	return gws
}

func (c *SIGCache) AttachedHTTPRoutes(gw *gatewayv1beta1.Gateway) []*gatewayv1beta1.HTTPRoute {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if gw == nil {
		return []*gatewayv1beta1.HTTPRoute{}
	}

	hrs := []*gatewayv1beta1.HTTPRoute{}
	for _, hr := range ActiveSIGs.HTTPRoute {
		for _, pr := range hr.Spec.ParentRefs {
			ns := hr.Namespace
			if pr.Namespace != nil {
				ns = string(*pr.Namespace)
			}
			if utils.Keyname(ns, string(pr.Name)) == utils.Keyname(gw.Namespace, gw.Name) {
				hrs = append(hrs, hr)
			}
		}
	}
	return hrs
}

func (c *SIGCache) GetRelatedObjs(gwObjs []*gatewayv1beta1.Gateway, hrObjs []*gatewayv1beta1.HTTPRoute, gwmap *map[string]*gatewayv1beta1.Gateway, hrmap *map[string]*gatewayv1beta1.HTTPRoute) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, gwObj := range gwObjs {
		if gwObj != nil {
			name := utils.Keyname(gwObj.Namespace, gwObj.Name)
			(*gwmap)[name] = c.GetGateway(name)
		}
		hrs := c.AttachedHTTPRoutes(gwObj)
		for _, hr := range hrs {
			name := utils.Keyname(hr.Namespace, hr.Name)
			if _, ok := (*hrmap)[name]; !ok {
				(*hrmap)[name] = c.GetHTTPRoute(name)
				c.GetRelatedObjs([]*gatewayv1beta1.Gateway{}, []*gatewayv1beta1.HTTPRoute{hr}, gwmap, hrmap)
			}
		}
	}

	for _, hrObj := range hrObjs {
		if hrObj != nil {
			name := utils.Keyname(hrObj.Namespace, hrObj.Name)
			(*hrmap)[name] = c.GetHTTPRoute(name)
		}
		gws := c.ParentRefsOf(hrObj)
		for _, gw := range gws {
			name := utils.Keyname(gw.Namespace, gw.Name)
			if _, ok := (*gwmap)[name]; !ok {
				(*gwmap)[name] = c.GetGateway(name)
				c.GetRelatedObjs([]*gatewayv1beta1.Gateway{gw}, []*gatewayv1beta1.HTTPRoute{}, gwmap, hrmap)
			}
		}
	}
}