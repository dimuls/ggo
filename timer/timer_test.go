package timer_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/someanon/ggo/timer"
)

var _ = Describe("Timer", func() {
	Describe("creating", func() {
		var err error
		Context("when base or byo-yomi are negative", func() {
			It("should be error", func() {
				_, err = NewTimer(Parameters{-1, 0, 0, 0}, Callbacks{})
				Expect(err).To(HaveOccurred())
				_, err = NewTimer(Parameters{1, -1, 1, 1}, Callbacks{})
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when both base and byo-yomi are zero", func() {
			It("should be error", func() {
				_, err = NewTimer(Parameters{0, 0, 1, 1}, Callbacks{})
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when not zero period or moves when zero byo-yomi", func() {
			It("should be error", func() {
				_, err = NewTimer(Parameters{1, 0, 1, 0}, Callbacks{})
				Expect(err).To(HaveOccurred())
				_, err = NewTimer(Parameters{1, 0, 0, 1}, Callbacks{})
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when zero or negative period when byo-yomi greater than zero", func() {
			It("should be error", func() {
				_, err = NewTimer(Parameters{0, 1, 0, 1}, Callbacks{})
				Expect(err).To(HaveOccurred())
				_, err = NewTimer(Parameters{0, 1, -1, 1}, Callbacks{})
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when zero or negative moves when byo-yomi greater than zero", func() {
			It("should be error", func() {
				_, err = NewTimer(Parameters{0, 1, 1, 0}, Callbacks{})
				Expect(err).To(HaveOccurred())
				_, err = NewTimer(Parameters{0, 1, 1, -1}, Callbacks{})
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when both period and moves are greater than one", func() {
			It("should be error", func() {
				_, err = NewTimer(Parameters{0, 1, 2, 2}, Callbacks{})
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when valid parameters", func() {
			var t *Timer
			It("should succeed", func() {
				t, err = NewTimer(Parameters{1, 0, 0, 0}, Callbacks{})
				Expect(err).ToNot(HaveOccurred())
				Expect(t).ToNot(BeNil())
				t, err = NewTimer(Parameters{0, 1, 1, 1}, Callbacks{})
				Expect(err).ToNot(HaveOccurred())
				Expect(t).ToNot(BeNil())
				t, err = NewTimer(Parameters{1, 1, 1, 1}, Callbacks{})
				Expect(err).ToNot(HaveOccurred())
				Expect(t).ToNot(BeNil())
				t, err = NewTimer(Parameters{0, 1, 10, 1}, Callbacks{})
				Expect(err).ToNot(HaveOccurred())
				Expect(t).ToNot(BeNil())
				t, err = NewTimer(Parameters{0, 1, 1, 10}, Callbacks{})
				Expect(err).ToNot(HaveOccurred())
				Expect(t).ToNot(BeNil())
			})
		})
	})
	Describe("running", func() {
		var (
			periodOver bool
			baseOver   bool
			over       bool
		)
		reset := func() {
			periodOver = false
			baseOver = false
			over = false
		}
		createTimer := func(base, byoYomi, periods, moves int) *Timer {
			t, err := NewTimer(Parameters{base, byoYomi, periods, moves}, Callbacks{
				OnPeriodOver: func() {
					periodOver = true
				},
				OnBaseOver: func() {
					baseOver = true
				},
				OnOver: func() {
					over = true
				},
			})
			Expect(err).ToNot(HaveOccurred())
			return t
		}
		BeforeEach(func() {
			reset()
		})
		Context("only base time", func() {
			It("works correct", func() {
				t := createTimer(1, 0, 0, 0)
				t.Switch()
				time.Sleep(500 * time.Millisecond)
				t.Switch()
				Expect(periodOver).To(BeFalse())
				Expect(baseOver).To(BeFalse())
				Expect(over).To(BeFalse())
				time.Sleep(600 * time.Millisecond)
				Expect(periodOver).To(BeFalse())
				Expect(baseOver).To(BeFalse())
				Expect(over).To(BeFalse())
				t.Switch()
				time.Sleep(499 * time.Millisecond)
				Expect(periodOver).To(BeFalse())
				Expect(baseOver).To(BeFalse())
				Expect(over).To(BeFalse())
				time.Sleep(1 * time.Millisecond)
				Expect(periodOver).To(BeFalse())
				Expect(baseOver).To(BeFalse())
				Expect(over).To(BeTrue())
			})
		})
		Context("base time with one byo-yomi period", func() {
			It("works correct", func() {
				t := createTimer(1, 1, 1, 1)
				t.Switch()
				time.Sleep(1001 * time.Millisecond)
				Expect(periodOver).To(BeFalse())
				Expect(baseOver).To(BeTrue())
				Expect(over).To(BeFalse())
				reset()
				time.Sleep(998 * time.Millisecond)
				t.Switch()
				Expect(periodOver).To(BeFalse())
				Expect(baseOver).To(BeFalse())
				Expect(over).To(BeFalse())
				t.Switch()
				time.Sleep(999 * time.Millisecond)
				Expect(periodOver).To(BeFalse())
				Expect(baseOver).To(BeFalse())
				Expect(over).To(BeFalse())
				reset()
				time.Sleep(2 * time.Millisecond)
				Expect(periodOver).To(BeFalse())
				Expect(baseOver).To(BeFalse())
				Expect(over).To(BeTrue())
			})
		})
		Context("only multiple byo-yomi periods", func() {
			It("works correct", func() {
				t := createTimer(0, 1, 2, 1)
				t.Switch()
				time.Sleep(1001 * time.Millisecond)
				Expect(periodOver).To(BeTrue())
				Expect(baseOver).To(BeFalse())
				Expect(over).To(BeFalse())
			})
		})
		Context("only canadian byo-yomi", func() {
			It("overs correct", func() {
				t := createTimer(0, 1, 1, 2)
				t.Switch()
				time.Sleep(500 * time.Millisecond)
				t.Switch()
				Expect(periodOver).To(BeFalse())
				Expect(baseOver).To(BeFalse())
				Expect(over).To(BeFalse())
				t.Switch()
				time.Sleep(501 * time.Millisecond)
				Expect(periodOver).To(BeFalse())
				Expect(baseOver).To(BeFalse())
				Expect(over).To(BeTrue())
			})
			It("works correct", func() {
				t := createTimer(0, 1, 1, 2)
				t.Switch()
				time.Sleep(500 * time.Millisecond)
				t.Switch()
				Expect(periodOver).To(BeFalse())
				Expect(baseOver).To(BeFalse())
				Expect(over).To(BeFalse())
				t.Switch()
				t.Switch()
				t.Switch()
				time.Sleep(999 * time.Millisecond)
				Expect(periodOver).To(BeFalse())
				Expect(baseOver).To(BeFalse())
				Expect(over).To(BeFalse())
			})
		})
	})
})
