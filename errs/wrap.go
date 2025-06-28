// Copyright © 2025 LegzDev.
//
// This file is part of BaitoMeBot (see https://github.com/legzdev/BaitoMeBot).
//
// BaitoMeBot is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// BaitoMeBot is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with BaitoMeBot. If not, see <https://www.gnu.org/licenses/>.

package errs

import "fmt"

type wrappedError struct {
	err error
	msg string
}

func (e *wrappedError) Error() string {
	return e.err.Error() + ": " + e.msg
}

func (e *wrappedError) Unwrap() error {
	return e.err
}

func Wrap(err error, v ...any) error {
	return &wrappedError{err: err, msg: fmt.Sprint(v...)}
}

func Wrapf(err error, format string, v ...any) error {
	return &wrappedError{err: err, msg: fmt.Sprintf(format, v...)}
}
