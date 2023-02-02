package main

	type hkzhaoImpl struct {}
	var _ hkzhao = (*hkzhaoImpl)(nil)
	
	func (impl *hkzhaoImpl) getResp (req, Req,ip int32) (resp Resp)  {
		return nil
	}

