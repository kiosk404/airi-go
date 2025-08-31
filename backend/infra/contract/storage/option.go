package storage

import (
	"time"
)

type GetOptFn func(option *GetOption)

type GetOption struct {
	Expire int64 //  seconds
}

func WithExpire(expire int64) GetOptFn {
	return func(o *GetOption) {
		o.Expire = expire
	}
}

type PutOption struct {
	ContentType        *string
	ContentEncoding    *string
	ContentDisposition *string
	ContentLanguage    *string
	Expires            *time.Time
	Tagging            map[string]string
	ObjectSize         int64
}

type PutOptFn func(option *PutOption)

func WithTagging(tag map[string]string) PutOptFn {
	return func(o *PutOption) {
		if len(tag) > 0 {
			o.Tagging = make(map[string]string, len(tag))
			for k, v := range tag {
				o.Tagging[k] = v
			}
		}
	}
}

func WithContentType(v string) PutOptFn {
	return func(o *PutOption) {
		o.ContentType = &v
	}
}

func WithObjectSize(v int64) PutOptFn {
	return func(o *PutOption) {
		o.ObjectSize = v
	}
}

func WithContentEncoding(v string) PutOptFn {
	return func(o *PutOption) {
		o.ContentEncoding = &v
	}
}

func WithContentDisposition(v string) PutOptFn {
	return func(o *PutOption) {
		o.ContentDisposition = &v
	}
}

func WithContentLanguage(v string) PutOptFn {
	return func(o *PutOption) {
		o.ContentLanguage = &v
	}
}

func WithExpires(v time.Time) PutOptFn {
	return func(o *PutOption) {
		o.Expires = &v
	}
}
