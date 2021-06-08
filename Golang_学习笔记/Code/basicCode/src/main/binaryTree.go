package main

import (
	"fmt"
	"tree"
)

/**
二叉树：中序遍历
 */

/**
遍历所有的节点的数值
 */
func MidTravel(tree *tree.Tree, ch chan int){
	if tree == nil{
		return
	}
	fmt.Println(tree.Val)
	ch <- tree.Val
	//左遍历
	MidTravel(tree.Left, ch)
	//右遍历
	MidTravel(tree.Right, ch)
}

/**
判断两棵树的节点值是否完全一致
 */
func IsSame(tree1 *tree.Tree, tree2 *tree.Tree) bool{
	if tree1 == nil && tree2 == nil{
		return true
	}
	if tree1 != nil && tree2 != nil && tree1.Val == tree2.Val{
		return IsSame(tree1.Left, tree2.Left) && IsSame(tree1.Right, tree2.Right)
	} else{
		return false
	}
}

func main() {
	//创建两颗二叉树
	tree1 := tree.Tree{Val:   1, Left:  nil, Right: nil}
	node1 := tree.Tree{Val:   2, Left:  nil, Right: nil}
	node2 := tree.Tree{Val:   3, Left:  nil, Right: nil}
	node3 := tree.Tree{Val:   4, Left:  nil, Right: nil}
	node4 := tree.Tree{Val:   5, Left:  nil, Right: nil}
	tree1.Left = &node1
	tree1.Right = &node2
	node1.Left = &node3
	node1.Right = &node4

	tree2 := tree.Tree{Val:   1, Left:  nil, Right: nil}
	node5 := tree.Tree{Val:   2, Left:  nil, Right: nil}
	node6 := tree.Tree{Val:   3, Left:  nil, Right: nil}
	node7 := tree.Tree{Val:   4, Left:  nil, Right: nil}
	node8 := tree.Tree{Val:   5, Left:  nil, Right: nil}
	tree2.Left = &node5
	tree2.Right = &node6
	node5.Left = &node7
	node5.Right = &node8

	//make(chan int)：默认不带缓冲区，即：只能存入一个数据，需要等待go协程取出后，才能解除阻塞状态
	ch1 := make(chan int, 5)
	ch2 := make(chan int, 5)
	fmt.Println("=====tree1====")
	MidTravel(&tree1, ch1)
	fmt.Println("=====tree2====")
	MidTravel(&tree2, ch2)

	fmt.Println("tree1 == tree2 ?  ", IsSame(&tree1, &tree2))

	fmt.Printf("%#v\n", tree1)
	fmt.Printf("%+v", tree1)

}