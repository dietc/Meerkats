package main

import (
	"testing"
	"s"
	"set"
)

func newPre() *s.Pre{
	str1 := "dir/doc1"
	str2 := "pic1.jpg"
    var pre *s.Pre = new(s.Pre)
	pre.Device = 0x01
	pre.P = make(map[string]*s.Unit)
	pre.P[str1] = new(s.Unit)
	pre.P[str2] = new(s.Unit)
	arr1 := [16]byte{62,238,247,95,134,12,98,97,1,23,50,153,132,75,115,133}
	pre.P[str1].Digest = arr1
	pre.P[str1].Idx = 1551931767918644362
	arr2 := [16]byte{141,221,57,244,109,128,93,240,169,147,223,23,244,177,183,2}
	pre.P[str2].Digest = arr2
	pre.P[str2].Idx = 1551931767923259098
	pre.Cmp = make(map[[16]byte]*set.Set)
	set1 := set.New()
	set1.Add(str1)
	pre.Cmp[arr1] = set1
	set2 := set.New()
	set2.Add(str2)
	pre.Cmp[arr2] = set2
	pre.Checked = make(map[string]bool)
	pre.Checked[str1] = false
	pre.Checked[str2] = false
	return pre
}



func TestCheckSpaceStatus(t *testing.T) {
	str1 := "dir/doc1"
	str2 := "pic1.jpg"
	arr1 := [16]byte{62,238,247,95,134,12,98,97,1,23,50,153,132,75,115,133}
	arr2 := [16]byte{141,221,57,244,109,128,93,240,169,147,223,23,244,177,183,2}
	
	//case0
	var space s.Space = make(map[string]*s.Unit)
	pret, newt := s.CheckSpaceStatus(newPre(), space)
	if pret[str1].Typ != 4 || pret[str2].Typ != 4 || len(newt) != 0 {
	    t.Error(pret[str1].Typ, pret[str2].Typ, len(newt))
	}
	
	//case1
	space[str1] = &s.Unit{arr1, 1551931767918644362}
	space[str2] = &s.Unit{arr2, 1551931767923259098}
	pret, newt = s.CheckSpaceStatus(newPre(), space)
	if pret[str1].Typ != 0 || pret[str2].Typ != 0 || len(newt) != 0{
		t.Error(pret[str1].Typ, pret[str2].Typ, len(newt))
	}
	
	//case3 
	delete(space, str1)
	space[str1] = &s.Unit{[16]byte{23,238,247,95,134,12,98,97,1,23,50,153,132,75,115,133}, 1551931767918644360}
    pret, newt = s.CheckSpaceStatus(newPre(), space)
	if pret[str1].Typ != 2 || pret[str2].Typ != 0 || len(newt) != 0{
		t.Error(pret[str1].Typ, pret[str2].Typ, len(newt))
	}
	
	//case4
	delete(space, str1)
	str3 := "dir/doc1new"
	space[str3] = &s.Unit{arr1, 1551931767918644365} 
	pret, newt = s.CheckSpaceStatus(newPre(), space)
	if pret[str1].Typ != 3 || pret[str2].Typ != 0 || len(newt) != 0{
		t.Error(pret[str1].Typ, pret[str2].Typ, len(newt))
	}
	
	//case5
	str4 := "extra"
	space[str4]  = &s.Unit{arr1, 155193176791244365}
	pret, newt = s.CheckSpaceStatus(newPre(), space)
	if pret[str1].Typ != 3 || pret[str2].Typ != 0 || len(newt) != 1{
		t.Error(pret[str1].Typ, pret[str2].Typ, newt[str4])
	}
}


//test cases of changed dir

func TestCheckCurrentStatus(t *testing.T) {
	str1 := "dir/doc1"
	str2 := "pic1.jpg"
	arr1 := [16]byte{62,238,247,95,134,12,98,97,1,23,50,153,132,75,115,133}
	arr2 := [16]byte{141,221,57,244,109,128,93,240,169,147,223,23,244,177,183,2}
	obj1 := s.FileObj{str1, 1, arr1}
	obj2 := s.FileObj{str2, 1, arr2}
	
	//case0
	var fileList map[string]s.FileObj = make(map[string]s.FileObj)
	pret, _ := s.CheckCurrentStatus(newPre(), fileList)
	if pret[str1].Typ != 4 || pret[str2].Typ !=4 {
		t.Error(pret[str1].Typ, pret[str2].Typ)
	}
	
	//case1	
	fileList[str1] = obj1
	fileList[str2] = obj2
	pret, _ = s.CheckCurrentStatus(newPre(), fileList)
	if pret[str1].Typ != 0 || pret[str2].Typ !=0 {
		t.Error(pret[str1].Typ, pret[str2].Typ)
	}
	
	//case2
	str3 := "extra"
	arr3 := [16]byte{51,238,247,95,134,12,98,97,1,23,50,153,132,75,115,123}
	obj3 := s.FileObj{str3, 1, arr3}
	fileList[str3] = obj3
	pret, newt := s.CheckCurrentStatus(newPre(), fileList)
	_,ok := newt[str3]
	if pret[str1].Typ != 0 || pret[str2].Typ !=0 || !ok {
		t.Error(pret[str1].Typ, pret[str2].Typ, newt[str3])
	}
	
	//case3
	delete(fileList, str2)
	arr4 := [16]byte{51,238,247,95,134,12,98,97,1,23,50,153,132,75,115,123}
	obj4 := s.FileObj{str2, 1, arr4}
	fileList[str2] = obj4
	pret, newt = s.CheckCurrentStatus(newPre(), fileList)
	_,ok = newt[str3]
	if pret[str1].Typ != 0 || pret[str2].Typ !=2 || !ok {
		t.Error(pret[str1].Typ, pret[str2].Typ, newt[str3])
	}
	
	//case4
	delete(fileList, str2)
	pret, newt = s.CheckCurrentStatus(newPre(), fileList)
	_,ok = newt[str3]
	if pret[str1].Typ != 0 || pret[str2].Typ !=4 || !ok {
		t.Error(pret[str1].Typ, pret[str2].Typ, newt[str3])
	}
	
	//case5
	obj5 := s.FileObj{str2, 1, obj1.Digest}
	fileList[str2] = obj5
	pret, newt = s.CheckCurrentStatus(newPre(), fileList)
	_,ok = newt[str3]
	if pret[str1].Typ != 0 || pret[str2].Typ !=2 || !ok   {
		t.Error(pret[str1].Typ, pret[str2].Typ, newt[str3])
	}
	
	//case6
	str4 := "dir/doc1new"
	obj6 := s.FileObj{str4, 1, obj1.Digest}
	delete(fileList, str1)
	fileList[str4] = obj6
	pret, newt = s.CheckCurrentStatus(newPre(), fileList)
	if pret[str1].Typ != 3 || pret[str2].Typ !=2 || !ok   {
		t.Error(pret[str1].Typ, pret[str2].Typ, newt[str3])
	}			
}



