#!/usr/bin/env python

from rospy_tutorials.srv import AddTwoInts
import rospy
import time


def main():
    rospy.init_node('add_two_ints_server')
    s = rospy.ServiceProxy('add_two_ints', AddTwoInts)
    print "Ready to add two ints."

    args = rospy.myargv()

    a = int(args[1])
    b = int(args[2])

    result = s(a, b)
    rospy.loginfo("{} + {} = {}".format(a, b, result.sum))
    time.sleep(1.0)

if __name__ == "__main__":
    main()

