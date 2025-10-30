import React, { useEffect, useState } from 'react'
import { useApi } from '../services/api'

const AdminDashboard = () => {
  const api = useApi()
  const [deposits, setDeposits] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let mounted = true
    const load = async () => {
      try {
        const res = await api.get('/admin/deposits/pending')
        if (mounted) setDeposits(res.data)
      } catch (err) {
        console.error(err)
      } finally {
        if (mounted) setLoading(false)
      }
    }
    load()
    return () => (mounted = false)
  }, [api])

  if (loading) return <div>Loading...</div>

  return (
    <div>
      <h2>Pending Deposits</h2>
      {deposits.length === 0 ? (
        <p>No pending deposits</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>User</th>
              <th>Total Points</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            {deposits.map((d) => (
              <tr key={d.ID || d.id}>
                <td>{d.ID || d.id}</td>
                <td>{d.user ? d.user.full_name : d.user_id}</td>
                <td>{d.total_points}</td>
                <td>{d.status}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}

export default AdminDashboard
