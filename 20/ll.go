package main

//
//type Node[T any] struct {
//  data T
//  prev *Node[T]
//  next *Node[T]
//}
//
//type LinkedList[T any] struct {
//  len  int
//  head *Node[T] // if head == nil ==> list empty
//  tail *Node[T]
//}
//
//func (t *LinkedList[T]) Add(data T) {
//  newNode := &Node[T]{
//    data: data,
//  }
//
//  if t.head == nil {
//    t.head = newNode
//    t.tail = newNode
//  } else {
//    lastTail := t.tail
//
//    lastTail.next = newNode
//    newNode.prev = lastTail
//    t.tail = newNode
//  }
//
//  t.len++
//}
//
//func (t *LinkedList[T]) ToSlice() []T {
//  ret := make([]T, t.len)
//
//  it := t.head
//  for it != nil {
//    ret = append(ret, it.data)
//    it = it.next
//  }
//
//  return ret
//}
//
//func (t *LinkedList[T]) Insert(pos int, data T) {
//  if t.len > pos {
//    panic("pos outside list")
//  }
//
//  prevNode := t.head
//  for i := 0; i < pos; i++ {
//    prevNode = prevNode.next
//  }
//
//  if prevNode == nil {
//    panic("Outside list")
//  }
//
//  newNode := Node[T]{
//    data: data,
//  }
//
//}
