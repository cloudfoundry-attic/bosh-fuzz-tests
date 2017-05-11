package parser_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EventLogParser", func() {
	Describe("ParseEventLog", func() {
		Context("When event log is invalid", func() {
			It("returns an error", func() {
				_, err := ParseEventLog("Garbage data")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when event log is valid", func() {
			Context("when the log has events", func() {
				It("should return the correct number of events", func() {
					events, err := ParseEventLog(EventLogWithEvents)
					Expect(err).ToNot(HaveOccurred())
					Expect(events).To(HaveLen(3))
				})
			})

			Context("when the log has no events", func() {
				It("should return the correct number of events", func() {
					events, err := ParseEventLog(EventLogWithNoEvents)
					Expect(err).ToNot(HaveOccurred())
					Expect(events).To(HaveLen(0))
				})
			})
		})
	})
	Describe("FindById", func() {
		events := Events{
			Event{Id: "a", ObjectName: "/1"},
			Event{Id: "b", ObjectName: "/2"},
			Event{Id: "c", ObjectName: "/3"},
		}
		Context("when Events has an Event with matching id", func() {
			It("should return the event", func() {
				event, err := events.FindById("/2")
				Expect(err).ToNot(HaveOccurred())
				Expect(event).ToNot(BeNil())
				Expect(event.Id).To(Equal("b"))
				Expect(event.ObjectName).To(Equal("/2"))
			})
		})
		Context("when Events do not have an Event with matching id", func() {
			It("should return an error", func() {
				_, err := events.FindById("999")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Event with ObjectId '/TestDirector/foo-deployment/999' not found"))
			})
		})
	})
})

var EventLogWithNoEvents = `
{
    "Tables": [
        {
            "Rows": []
        }
    ]
}`

var EventLogWithEvents = `
{
    "Tables": [
        {
            "Rows": [
               {
                    "action": "create",
                    "context": "id: \"13\"\nname: /TestDirector/foo-deployment/nyAtiEiFR383GgdIXqjU",
                    "deployment": "foo-deployment",
                    "error": "",
                    "id": "65",
                    "instance": "",
                    "object_name": "/TestDirector/foo-deployment/nyAtiEiFR383GgdIXqjU",
                    "object_type": "variable",
                    "task_id": "5",
                    "time": "Mon May  8 15:17:23 UTC 2017",
                    "user": "test"
                },
                {
                    "action": "acquire",
                    "context": "",
                    "deployment": "foo-deployment",
                    "error": "",
                    "id": "61",
                    "instance": "",
                    "object_name": "lock:deployment:foo-deployment",
                    "object_type": "lock",
                    "task_id": "5",
                    "time": "Mon May  8 15:17:23 UTC 2017",
                    "user": "test"
                },
                {
                    "action": "update",
                    "context": "",
                    "deployment": "foo-deployment",
                    "error": "",
                    "id": "60",
                    "instance": "",
                    "object_name": "foo-deployment",
                    "object_type": "deployment",
                    "task_id": "5",
                    "time": "Mon May  8 15:17:23 UTC 2017",
                    "user": "test"
                }
            ]
        }
    ]
}`
