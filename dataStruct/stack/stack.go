package stack

import "sync"

/*
	stack栈结构体
	包含动态数组和该数组的顶部指针
	顶部指针指向实际顶部元素的下一位置
	当删除节点时仅仅需要下移顶部指针一位即可
	当新增结点时优先利用冗余空间
	当冗余空间不足时先倍增空间至2^16，超过后每次增加2^16的空间
	删除结点后如果冗余超过2^16,则释放掉
	删除后若冗余量超过使用量，也释放掉冗余空间
 */

type Stack struct {
	data  []interface{} //用于存储元素的动态数组
	top   uint64        //顶部指针
	cap   uint64        //动态数组的实际空间
	needSync bool       // 是否需要并发控制
	mutex sync.Mutex    //并发控制锁
}

func NewStack() (s *Stack) {
	return &Stack{
		data:  make([]interface{}, 1, 1),
		top:   0,
		cap:   1,
		needSync: false,
	}
}

func NewSynStack() (s *Stack) {
	return &Stack{
		data:  make([]interface{}, 1, 1),
		top:   0,
		cap:   1,
		needSync: true,
		mutex: sync.Mutex{},
	}
}

func (s *Stack) Size() (num uint64) {
	if s == nil {
		return 0
	}
	return s.top
}

func (s *Stack) Pop() (e interface{},ok bool) {
	if s == nil {
		return nil,false
	}
	if s.Size() == 0 {
		return nil,false
	}

	if s.needSync {
		s.mutex.Lock()
		defer s.mutex.Unlock()
	}

	e = s.data[s.top-1]
	s.top--

	if s.cap-s.top >= 65536 {
		//容量和实际使用差值超过2^16时,容量直接减去2^16
		s.cap -= 65536
		tmp := make([]interface{}, s.cap, s.cap)
		copy(tmp, s.data)
		s.data = tmp
	} else if s.top*2 < s.cap {
		//实际使用长度是容量的一半时,进行折半缩容
		s.cap /= 2
		tmp := make([]interface{}, s.cap, s.cap)
		copy(tmp, s.data)
		s.data = tmp
	}

	return e,true
}

func (s *Stack) Push(e interface{}) {
	if s == nil {
		s = NewStack()
	}

	if s.needSync {
		s.mutex.Lock()
		defer s.mutex.Unlock()
	}

	if s.top < s.cap {
		//还有冗余,直接添加
		s.data[s.top] = e
	} else {
		//冗余不足,需要扩容
		if s.cap <= 65536 {
			//容量翻倍
			if s.cap == 0 {
				s.cap = 1
			}
			s.cap *= 2
		} else {
			//容量增加2^16
			s.cap += 65536
		}
		//复制扩容前的元素
		tmp := make([]interface{}, s.cap, s.cap)
		copy(tmp, s.data)
		s.data = tmp
		s.data[s.top] = e
	}
	s.top++
}

