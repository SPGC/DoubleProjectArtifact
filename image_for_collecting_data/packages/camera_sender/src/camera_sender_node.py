#!/usr/bin/env python3

import os

import rospy
from duckietown.dtros import DTROS, NodeType, TopicType, DTParam, ParamType
from sensor_msgs.msg import CompressedImage

import requests

import cv2
from cv_bridge import CvBridge


class CameraSenderNode(DTROS):

    def __init__(self, node_name):
        super(CameraSenderNode, self).__init__(node_name=node_name, node_type=NodeType.VISUALIZATION)
        rospy.get_param("~ip", "192.168.73.254:8080") # To initialize paramter
        self._vehicle_name = os.environ['VEHICLE_NAME']
        self._camera_topic = f"/{self._vehicle_name}/camera_node/image/compressed"
        self._bridge = CvBridge()
        self.sub = rospy.Subscriber(self._camera_topic, CompressedImage, self.callback, buff_size=10000000,
                                    queue_size=1)
        self.ip = DTParam("~ip", param_type=ParamType.STRING)

    def callback(self, msg):
        print("Callback called, ip = ", self.ip.value)
        img = cv2.cvtColor(self._bridge.compressed_imgmsg_to_cv2(msg), cv2.COLOR_BGR2RGB)
        url = "http://" + self.ip.value
        array_to_send = img.tolist()
        data = {"type": "img", "data": array_to_send}
        try:
            response = requests.post(url, json=data)
        except:
            print("Something went went wrong while sending")
            return
        if not response.status_code == 200:
            print("Something went wrong, status code: " + str(response.status_code))
        else:
            print("Data sent")


if __name__ == '__main__':
    node = CameraSenderNode(node_name='camera_sender_node')
    rospy.spin()
